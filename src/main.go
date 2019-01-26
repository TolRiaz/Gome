package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

)

type Mode int

var(
	charImage		*assets.Image
)

type Game struct {
	mode Mode

	// The Character's position
	x int
	y int

	// Camera
	camX int
	camY int

	gameoverCounts int
}

func NewGame() *Game  {
	var g Game
	g.init()
	return g
}

func (g *Game) init() {
	g.x = 0
	g.y = 0
	// g.camx = 0
	// g.camy = 0
}

/*
func (g *Game) drawChar(screen *assets.Image) {
	op := &assets.DrawImageOptions{}
	w, h := charImage.Size()
	op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	op.GeoM.Rotate(float64(g.vy16) / 96.0 * math.Pi / 6)
	op.GeoM.Translate(float64(w)/2.0, float64(h)/2.0)
	op.GeoM.Translate(float64(g.x16/16.0)-float64(g.cameraX), float64(g.y16/16.0)-float64(g.cameraY))
	op.Filter = assets.FilterLinear
	screen.DrawImage(charImage, op)
}
*/



