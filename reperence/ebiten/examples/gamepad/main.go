// Copyright 2015 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build example jsgo

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	gamepadIDs = map[int]struct{}{}
)

func update(screen *ebiten.Image) error {
	// Log the gamepad connection events.
	for _, id := range inpututil.JustConnectedGamepadIDs() {
		log.Printf("gamepad connected: id: %d", id)
		gamepadIDs[id] = struct{}{}
	}
	for id := range gamepadIDs {
		if inpututil.IsGamepadJustDisconnected(id) {
			log.Printf("gamepad disconnected: id: %d", id)
			delete(gamepadIDs, id)
		}
	}

	ids := ebiten.GamepadIDs()
	axes := map[int][]string{}
	pressedButtons := map[int][]string{}

	for _, id := range ids {
		maxAxis := ebiten.GamepadAxisNum(id)
		for a := 0; a < maxAxis; a++ {
			v := ebiten.GamepadAxis(id, a)
			axes[id] = append(axes[id], fmt.Sprintf("%d:%0.2f", a, v))
		}
		maxButton := ebiten.GamepadButton(ebiten.GamepadButtonNum(id))
		for b := ebiten.GamepadButton(id); b < maxButton; b++ {
			if ebiten.IsGamepadButtonPressed(id, b) {
				pressedButtons[id] = append(pressedButtons[id], strconv.Itoa(int(b)))
			}

			// Log button events.
			if inpututil.IsGamepadButtonJustPressed(id, b) {
				log.Printf("button pressed: id: %d, button: %d", id, b)
			}
			if inpututil.IsGamepadButtonJustReleased(id, b) {
				log.Printf("button released: id: %d, button: %d", id, b)
			}
		}
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Draw the current gamepad status.
	str := ""
	if len(ids) > 0 {
		for _, id := range ids {
			str += fmt.Sprintf("Gamepad (ID: %d):\n", id)
			str += fmt.Sprintf("  Axes:    %s\n", strings.Join(axes[id], ", "))
			str += fmt.Sprintf("  Buttons: %s\n", strings.Join(pressedButtons[id], ", "))
			str += "\n"
		}
	} else {
		str = "Please connect your gamepad."
	}
	ebitenutil.DebugPrint(screen, str)
	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Gamepad (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}
