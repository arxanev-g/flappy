package main

type Bird struct {
	X, Y   float64
	VY     float64
	Width  float64
	Height float64
}

func NewBird(screenWidth, screenHeight float64) *Bird {
	return &Bird{
		X:      screenWidth * 0.1,
		Y:      screenHeight / 2,
		VY:     0,
		Width:  50,
		Height: 50,
	}
}

func (b *Bird) Update(gravity, delta float64) {
	b.VY += gravity * delta
	b.Y += b.VY
}

func (b *Bird) Jump(jumpForce float64) {
	b.VY = jumpForce
}
