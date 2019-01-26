package main

import(
	"image"
	"log"
	"math"
	"math/rand"

	_ "image/png"
)

type Game struct {
	character struct {
	x		float32		// x-offset
	y		float32		// y-offset
	v		float32		// velocity
	dead	bool		// Dead or Alive
	// deadTime clock.Time // when the character died
	}
}

// struct clock {
//		type Time int32
// }

func NewGame() *Game {
	var g Game
	g.reset()
	return &g
}

func (g *Game) reset() {
	g.gopher.x = 0
	g.gopher.y = 0
	g.gopher.v = 0

	g.gopher.dead = false
}

