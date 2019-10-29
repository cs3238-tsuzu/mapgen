package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	"golang.org/x/image/bmp"
)

type GroundKind = int

const (
	Ground GroundKind = iota
	Grass
)

func calcGroundKind(col color.Color) GroundKind {
	r, g, b, a := get8bitColor(col)

	if isNormalGround(r, g, b, a) {
		return Ground
	}

	if isRespawnPosition(r, g, b, a) {
		return Ground
	}

	if isWalkableGrass(r, g, b, a) {
		return Grass
	}

	if isNormalGrass(r, g, b, a) {
		return Grass
	}

	panic("unknown color: " + fmt.Sprint(r, g, b, a))
}

func isWall(img image.Image, y, x int) bool {
	col := img.At(x, y)
	r, g, b, a := get8bitColor(col)

	if isNormalGrass(r, g, b, a) {
		return true
	}

	return false
}

func isObject(img image.Image, y, x int) bool {
	if !isWall(img, y, x) {
		return false
	}

	if isWall(img, y+1, x) ||
		isWall(img, y-1, x) ||
		isWall(img, y, x+1) ||
		isWall(img, y, x-1) {
		return true
	}

	return false
}

func getGroundKind(img image.Image, y, x int) GroundKind {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	if width := img.Bounds().Dx(); width <= x {
		x = width - 1
	}

	if height := img.Bounds().Dy(); height <= y {
		y = height - 1
	}

	return calcGroundKind(img.At(x, y))
}

func calcFilename(img image.Image, y, x int) string {
	switch getGroundKind(img, y, x) {
	case Grass:
		return "3333"
	}

	bs := make([]int, 4)
	for i := 0; i < 4; i++ {
		ay := int(-math.Cos(math.Pi / 2 * float64(i)))
		ax := int(math.Sin(math.Pi / 2 * float64(i)))

		by := int(-math.Sqrt(2) * math.Cos(math.Pi/4+math.Pi/2*float64(i)))
		bx := int(math.Sqrt(2) * math.Sin(math.Pi/4+math.Pi/2*float64(i)))

		cy := int(-math.Cos(math.Pi / 2 * float64(i+1)))
		cx := int(math.Sin(math.Pi / 2 * float64(i+1)))

		kind := getGroundKind(img, y+ay, x+ax) | getGroundKind(img, y+by, x+bx) | getGroundKind(img, y+cy, x+cx)

		bs[i] |= kind << 1
		bs[(i+1)%4] |= kind
	}

	name := ""
	for i := range bs {
		name = name + strconv.Itoa(bs[i])
	}

	return name
}

func createTileMap(base string, img image.Image) [][]TextureName {
	ret := [][]TextureName{}

	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	for y := 0; y < height; y++ {

		r := []TextureName{}
		for x := 0; x < width; x++ {
			r = append(r, TextureName(path.Join(base, calcFilename(img, y, x)+".png")))
		}

		ret = append(ret, r)
	}

	return ret
}

type ObjectJSON struct {
	Objects     map[string]string `json:"objects"`
	TileBase    string            `json:"tile_base"`
	Transparent string            `json:"transparent"`
}

func colorString(col color.Color) string {
	r, g, b, a := get8bitColor(col)

	// fmt.Println(r, g, b, a)
	if a == 0 {
		return ""
	}

	return strings.ToUpper(fmt.Sprintf("%02x%02x%02x", r, g, b))
}

func createObjectMap(oj *ObjectJSON, grdimg, objimg image.Image) [][]TextureNameOrNullString {
	ret := [][]TextureNameOrNullString{}

	width, height := grdimg.Bounds().Dx(), grdimg.Bounds().Dy()

	for y := 0; y < height; y++ {
		r := []TextureNameOrNullString{}
		for x := 0; x < width; x++ {
			if isObject(grdimg, y, x) {
				r = append(r, TextureNameOrNullString(oj.Transparent))
			} else {
				r = append(r, "")
			}
		}

		ret = append(ret, r)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			col := colorString(objimg.At(x, y))

			if col == "" {
				continue
			}

			obj, found := oj.Objects[col]

			if !found {
				panic("unknown color: " + col + fmt.Sprint(" for ", y, x))
			}

			ret[y][x] = TextureNameOrNullString(obj)
		}
	}

	return ret
}

func getRespawnPositions(grdimg image.Image) []RespawnPosition {
	width, height := grdimg.Bounds().Dx(), grdimg.Bounds().Dy()

	poss := []RespawnPosition{}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if isRespawnPosition(get8bitColor(grdimg.At(x, y))) {
				poss = append(poss, RespawnPosition{
					X: x,
					Y: y,
				})
			}
		}
	}

	return poss
}

func openTemplateJSON() *MapJSON {
	fp, err := os.Open("template.json")

	if err != nil {
		panic(err)
	}

	var mj MapJSON
	if err := json.NewDecoder(fp).Decode(&mj); err != nil {
		panic(err)
	}

	return &mj
}

func openObjectsJSON() *ObjectJSON {
	fp, err := os.Open("objects.json")

	if err != nil {
		panic(err)
	}

	var mj ObjectJSON
	if err := json.NewDecoder(fp).Decode(&mj); err != nil {
		panic(err)
	}

	return &mj
}

var tiles = []string{
	"0000.png",
	"0021.png",
	"0210.png",
	"0231.png",
	"1002.png",
	"1023.png",
	"1212.png",
	"1233.png",
	"2100.png",
	"2121.png",
	"2310.png",
	"2331.png",
	"3102.png",
	"3123.png",
	"3312.png",
	"3333.png",
}

func openImage(name string) image.Image {
	fp, err := os.Open(name)

	if err != nil {
		panic(err)
	}

	img, err := bmp.Decode(fp)

	if err != nil {
		panic(err)
	}

	defer fp.Close()

	return img
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println(os.Args[0], "[path to map bmp]", "[path to objects bmp]")

		os.Exit(1)
	}

	output := openTemplateJSON()
	oj := openObjectsJSON()

	grd := os.Args[1]
	obj := os.Args[2]

	grdimg := openImage(grd)
	objimg := openImage(obj)

	/*width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	const cellSize = 30

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fmt.Println(filepath.Join("magical", calcFilename(img, y, x)))

		}
	}*/

	output.TextureNameToURL[TextureName(oj.Transparent)] = oj.Transparent

	for i := range tiles {
		output.TextureNameToURL[TextureName(path.Join(oj.TileBase, tiles[i]))] = path.Join(oj.TileBase, tiles[i])
	}

	for _, p := range oj.Objects {
		output.TextureNameToURL[TextureName(p)] = p
	}

	output.Layers.TileMap = createTileMap(oj.TileBase, grdimg)
	output.Layers.Object = createObjectMap(oj, grdimg, objimg)

	output.RespawnPositions = getRespawnPositions(grdimg)

	out, err := os.Create("map.json")

	if err != nil {
		panic(err)
	}
	defer out.Close()

	if err := json.NewEncoder(out).Encode(output); err != nil {
		panic(err)
	}
}
