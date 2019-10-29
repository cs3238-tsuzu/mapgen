package main

import "image/color"

func get8bitColor(col color.Color) (r, g, b, a uint32) {
	r, g, b, a = col.RGBA()
	r /= 256
	g /= 256
	b /= 256
	a /= 256

	return
}
