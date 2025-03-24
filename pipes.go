package main

import (
	"math/rand"
	"syscall/js"
)

type Pipe struct {
	X         float64
	Top       float64
	Bottom    float64
	Width     float64
	Speed     float64
	GapHeight float64
	PipeColor string
}

func NewPipe(screenWidth, screenHeight, speed float64) *Pipe {
	gap := screenHeight * 0.25
	minHeight := screenHeight * 0.15
	top := rand.Float64()*(screenHeight-gap-minHeight) + minHeight

	return &Pipe{
		X:         screenWidth,
		Top:       top,
		Bottom:    top + gap,
		Width:     screenWidth * 0.05,
		Speed:     speed,
		GapHeight: gap,
		PipeColor: "#2ecc71",
	}
}

func (p *Pipe) Update() {
	p.X -= p.Speed
}

func (p *Pipe) IsOffscreen() bool {
	return p.X < -p.Width
}

func (p *Pipe) CheckCollision(bird *Bird) bool {
	overlapX := bird.X < p.X+p.Width &&
		bird.X+bird.Width > p.X

	overlapTop := bird.Y < p.Top
	overlapBottom := bird.Y+bird.Height > p.Bottom

	return overlapX && (overlapTop || overlapBottom)
}

func (p *Pipe) Draw(ctx js.Value) {
	ctx.Set("fillStyle", p.PipeColor)
	ctx.Call("fillRect", p.X, 0, p.Width, p.Top)
	ctx.Call("fillRect", p.X, p.Bottom, p.Width, height-p.Bottom)
}
