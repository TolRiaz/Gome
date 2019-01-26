// Copyright 2016 Hajime Hoshi
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

// This is an example to implement an audio player.
// See examples/wav for a simpler example to play a sound file.

package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	raudio "github.com/hajimehoshi/ebiten/examples/resources/audio"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240

	// This sample rate doesn't match with wav/vorbis's sample rate,
	// but decoders adjust them.
	sampleRate = 48000
)

var (
	playerBarColor     = color.RGBA{0x80, 0x80, 0x80, 0xff}
	playerCurrentColor = color.RGBA{0xff, 0xff, 0xff, 0xff}
)

// Player represents the current audio state.
type Player struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte
	volume128    int
}

func playerBarRect() (x, y, w, h int) {
	w, h = 300, 4
	x = (screenWidth - w) / 2
	y = screenHeight - h - 16
	return
}

func NewPlayer(audioContext *audio.Context) (*Player, error) {
	const bytesPerSample = 4 // TODO: This should be defined in audio package
	s, err := vorbis.Decode(audioContext, audio.BytesReadSeekCloser(raudio.Ragtime_ogg))
	if err != nil {
		return nil, err
	}
	p, err := audio.NewPlayer(audioContext, s)
	if err != nil {
		return nil, err
	}
	player := &Player{
		audioContext: audioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / sampleRate,
		volume128:    128,
		seCh:         make(chan []byte),
	}
	if player.total == 0 {
		player.total = 1
	}
	player.audioPlayer.Play()
	go func() {
		s, err := wav.Decode(audioContext, audio.BytesReadSeekCloser(raudio.Jab_wav))
		if err != nil {
			log.Fatal(err)
			return
		}
		b, err := ioutil.ReadAll(s)
		if err != nil {
			log.Fatal(err)
			return
		}
		player.seCh <- b
	}()
	return player, nil
}

func (p *Player) update() error {
	select {
	case p.seBytes = <-p.seCh:
		close(p.seCh)
		p.seCh = nil
	default:
	}

	if p.audioPlayer.IsPlaying() {
		p.current = p.audioPlayer.Current()
	}
	p.seekBarIfNeeded()
	p.switchPlayStateIfNeeded()
	p.playSEIfNeeded()
	p.updateVolumeIfNeeded()

	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		b := ebiten.IsRunnableInBackground()
		ebiten.SetRunnableInBackground(!b)
	}
	return nil
}

func (p *Player) playSEIfNeeded() {
	if p.seBytes == nil {
		// Bytes for the SE is not loaded yet.
		return
	}

	if !inpututil.IsKeyJustPressed(ebiten.KeyP) {
		return
	}
	sePlayer, _ := audio.NewPlayerFromBytes(p.audioContext, p.seBytes)
	sePlayer.Play()
}

func (p *Player) updateVolumeIfNeeded() {
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		p.volume128--
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		p.volume128++
	}
	if p.volume128 < 0 {
		p.volume128 = 0
	}
	if 128 < p.volume128 {
		p.volume128 = 128
	}
	p.audioPlayer.SetVolume(float64(p.volume128) / 128)
}

func (p *Player) switchPlayStateIfNeeded() {
	if !inpututil.IsKeyJustPressed(ebiten.KeyS) {
		return
	}
	if p.audioPlayer.IsPlaying() {
		p.audioPlayer.Pause()
		return
	}
	p.audioPlayer.Play()
}

func (p *Player) seekBarIfNeeded() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}

	// Calculate the next seeking position from the current cursor position.
	x, y := ebiten.CursorPosition()
	bx, by, bw, bh := playerBarRect()
	const padding = 4
	if y < by-padding || by+bh+padding <= y {
		return
	}
	if x < bx || bx+bw <= x {
		return
	}
	pos := time.Duration(x-bx) * p.total / time.Duration(bw)
	p.current = pos
	p.audioPlayer.Seek(pos)
}

func (p *Player) close() error {
	return p.audioPlayer.Close()
}

func (p *Player) draw(screen *ebiten.Image) {
	// Draw the bar.
	x, y, w, h := playerBarRect()
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), playerBarColor)

	// Draw the cursor on the bar.
	c := p.current
	cw, ch := 4, 10
	cx := int(time.Duration(w)*c/p.total) + x - cw/2
	cy := y - (ch-h)/2
	ebitenutil.DrawRect(screen, float64(cx), float64(cy), float64(cw), float64(ch), playerCurrentColor)

	// Compose the curren time text.
	m := (c / time.Minute) % 100
	s := (c / time.Second) % 60
	currentTimeStr := fmt.Sprintf("%02d:%02d", m, s)

	// Draw the debug message.
	msg := fmt.Sprintf(`FPS: %0.2f
Press S to toggle Play/Pause
Press P to play SE
Press Z or X to change volume of the music
Press B to switch the run-in-background state
Current Time: %s`, ebiten.CurrentFPS(), currentTimeStr)
	ebitenutil.DebugPrint(screen, msg)
}

var (
	musicPlayer *Player
)

func update(screen *ebiten.Image) error {
	if err := musicPlayer.update(); err != nil {
		return err
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	musicPlayer.draw(screen)
	return nil
}

func main() {
	audioContext, err := audio.NewContext(sampleRate)
	if err != nil {
		log.Fatal(err)
	}

	musicPlayer, err = NewPlayer(audioContext)
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Audio (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
	if err := musicPlayer.close(); err != nil {
		log.Fatal(err)
	}
}
