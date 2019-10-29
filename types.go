package main

import (
	"database/sql/driver"
	"encoding/json"
)

// このクラスにはマップ状態を表すMapJSONを定義しています。
// TypeScript, Go, Rust で同じ型を定義しています。

// TextureName はテクスチャ名を表す型です
type TextureName string

// TextureNameOrNullString は TextureName | "" を表す型です
type TextureNameOrNullString string

// AnimationName はアニメーション名を表す型です
type AnimationName string

// EnemyName は敵の名前を表す型です
type EnemyName string

// TextureNameOrAnimationName は TextureName | AnimationName を表す型です
type TextureNameOrAnimationName string

// AnimationName4 は4方向のアニメーションです
type AnimationName4 struct {
	Front AnimationName `json:"front"`
	Left  AnimationName `json:"left"`
	Right AnimationName `json:"right"`
	Back  AnimationName `json:"back"`
}

// AnimationJSON は1つのアニメーションに関する情報を表す構造体です
type AnimationJSON struct {
	Textures   []TextureName `json:"textures"`
	Width      float32       `json:"width"`
	Height     float32       `json:"height"`
	IntervalMs float32       `json:"interval_ms"`
}

// EnemyDataJSON は敵の種類を表す構造体です
type EnemyDataJSON struct {
	AnimationFront AnimationName `json:"animation_front"`
	AnimationLeft  AnimationName `json:"animation_left"`
	AnimationRight AnimationName `json:"animation_right"`
	AnimationBack  AnimationName `json:"animation_back"`
}

// RectCollisionJSON は矩形の当たり判定を表す構造体です
type RectCollisionJSON struct {
	Width   float32 `json:"width"`
	Height  float32 `json:"height"`
	CenterX float32 `json:"center_x"`
	CenterY float32 `json:"center_y"`
}

// CircleCollisionJSON は円形の当たり判定を表す構造体です
type CircleCollisionJSON struct {
	Radius  float32 `json:"radius"`
	CenterX float32 `json:"center_x"`
	CenterY float32 `json:"center_y"`
}

// CollisionsJSON は1つのオブジェクトの当たり判定を表す構造体です。
// 複数の当たり判定を合成できます。
type CollisionsJSON struct {
	Rect   []RectCollisionJSON   `json:"rect"`
	Circle []CircleCollisionJSON `json:"circle"`
}

// AppearanceJSON はマップ上でのプレイヤーの見た目を表す構造体でしゅ
type AppearanceJSON struct {
	HairType   string `json:"hair_type"`
	HairColor  string `json:"hair_color"`
	EyeType    string `json:"eye_type"`
	EyeColor   string `json:"eye_color"`
	Cloth      string `json:"cloth"`
	ClothColor string `json:"cloth_color"`
}

// UserJSON はマップ上のユーザーを表す構造体です
type UserJSON struct {
	UID  int    `json:"uid"`
	Name string `json:"name"`

	Time       int64           `json:"time"`
	PositionX  float32         `json:"position_x"`
	PositionY  float32         `json:"position_y"`
	VelocityX  float32         `json:"velocity_x"`
	VelocityY  float32         `json:"velocity_y"`
	Mass       float32         `json:"mass"`
	Appearance *AppearanceJSON `json:"appearance"`

	MaxHp      int      `json:"max_hp"`
	Hp         int      `json:"hp"`
	Attack     int      `json:"attack"`
	Defence    int      `json:"defence"`
	Experience int      `json:"experience"`
	Level      int      `json:"level"`
	Chat       string   `json:"chat"`
	Items      []string `json:"items"`

	NowLoading bool `json:"now_loading"`

	ActiveNPCID int    `json:"-"`
	NPCMode     string `json:"-"`
}

func (us *UserJSON) Scan(value interface{}) error {
	b := []byte{}
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	}

	return json.Unmarshal(b, us)
}

func (us *UserJSON) Value() (driver.Value, error) {
	return json.Marshal(us)
}

// NPCJSON はマップ上のNPCを表す構造体です
type NPCJSON struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	Time      int64          `json:"time"`
	PositionX float32        `json:"position_x"`
	PositionY float32        `json:"position_y"`
	VelocityX float32        `json:"velocity_x"`
	VelocityY float32        `json:"velocity_y"`
	Mass      float32        `json:"mass"`
	Animation AnimationName4 `json:"animation"`
}

// EnemyJSON はマップ上の敵を表す構造体です
type EnemyJSON struct {
	ID   int       `json:"id"`
	Name EnemyName `json:"name"`

	Time      int64   `json:"time"`
	PositionX float32 `json:"position_x"`
	PositionY float32 `json:"position_y"`
	VelocityX float32 `json:"velocity_x"`
	VelocityY float32 `json:"velocity_y"`
	Hp        int     `json:"hp"`
}

// CharacterJSON はマップ上のキャラクタの一覧を保持する構造体です
type CharacterJSON struct {
	User  map[int]*UserJSON  `json:"user"`
	NPC   map[int]*NPCJSON   `json:"npc"`
	Enemy map[int]*EnemyJSON `json:"enemy"`
}

// LayersJSON はマップの各レイヤのデータをまとめる構造体です
type LayersJSON struct {
	TileMap   [][]TextureName             `json:"tile_map"`  // mutable
	Object    [][]TextureNameOrNullString `json:"object"`    // mutable
	Character CharacterJSON               `json:"character"` // mutable
}

type RespawnPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// MapJSON で1つのマップの状態を表します
type MapJSON struct {
	YasunaDataJSON       *YasunaDataJSON                                `json:"yasuna_data_json"`       // immutable
	CharacterResourceURL string                                         `json:"character_resource_url"` // immutable
	TextureNameToURL     map[TextureName]string                         `json:"texture_name_to_url"`    // immutable
	TextureScale         map[TextureName]float32                        `json:"texture_scale"`          // immutable
	Animations           map[AnimationName]*AnimationJSON               `json:"animations"`             // immutable
	EnemyData            map[EnemyName]*EnemyDataJSON                   `json:"enemy_data"`
	Collisions           map[TextureNameOrAnimationName]*CollisionsJSON `json:"collisions"` // immutable
	Layers               LayersJSON                                     `json:"layers"`     // mutable
	RespawnPositions     []RespawnPosition                              `json:"respawn_positions"`
}

// YasunaDataJSON は tools/character-gen から生成されたyasuna.jsonの型です
type YasunaDataJSON struct {
	TextureName interface{} `json:"texture_name"`
	ImageSize   []int       `json:"image_size"`
	Layer       struct {
		LayerCount int `json:"layer_count"`
		Layers     []struct {
			Offset      []int  `json:"offset"`
			Clip        []int  `json:"clip"`
			TextureName string `json:"texture_name"`
		} `json:"layers"`
	} `json:"layer"`
	States struct {
		Tagged map[string]struct {
			TotalLayers int   `json:"total_layers"`
			TotalModes  int   `json:"total_modes"`
			AllLayers   []int `json:"all_layers"`
			Modes       map[string]struct {
				TotalLayers int   `json:"total_layers"`
				Layers      []int `json:"layers"`
			} `json:"modes"`
		} `json:"tagged"`
		Unique map[string]struct {
			TotalLayers int   `json:"total_layers"`
			Layers      []int `json:"layers"`
		} `json:"unique"`
		Animation map[string]struct {
			TotalLayers int      `json:"total_layers"`
			TotalFrames int      `json:"total_frames"`
			Layers      []int    `json:"layers"`
			Frames      [][]bool `json:"frames"`
		} `json:"animation"`
	} `json:"states"`
}

// NewMapJSON は空のMapJSONを生成する関数です
func NewMapJSON() *MapJSON {

	return &MapJSON{
		TextureNameToURL: make(map[TextureName]string),
		TextureScale:     make(map[TextureName]float32),
		Animations:       make(map[AnimationName]*AnimationJSON),
		EnemyData:        make(map[EnemyName]*EnemyDataJSON),
		Collisions:       make(map[TextureNameOrAnimationName]*CollisionsJSON),
		Layers: LayersJSON{
			TileMap: make([][]TextureName, 0, 2500),
			Object:  make([][]TextureNameOrNullString, 0, 2500),
			Character: CharacterJSON{
				User:  make(map[int]*UserJSON),
				NPC:   make(map[int]*NPCJSON),
				Enemy: make(map[int]*EnemyJSON),
			},
		},
	}
}
