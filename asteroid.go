package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Asteroid struct {
	Position Vector2
	Velocity Vector2
	Rotation float64
	Size     float64
}

func NewAsteroid(pos Vector2, vel Vector2, size float64) *Asteroid {
	return &Asteroid{
		Position: pos,
		Velocity: vel,
		Rotation: 0,
		Size:     size,
	}
}

func (a *Asteroid) Update() {
	a.Position = a.Position.Add(a.Velocity)
	a.Rotation += AsteroidRotationSpeed

	// if a.Position.X < -a.Size {
	// 	a.Position.X += ScreenWidth + a.Size*2
	// }
	// if a.Position.X > ScreenWidth+a.Size {
	// 	a.Position.X -= ScreenWidth + a.Size*2
	// }
	// if a.Position.Y < -a.Size {
	// 	a.Position.Y += ScreenHeight + a.Size*2
	// }
	// if a.Position.Y > -a.Size {
	// 	a.Position.Y -= ScreenHeight + a.Size*2
	// }
}

func (a *Asteroid) Draw(screen *ebiten.Image) {
	vector.FillCircle(
		screen,
		float32(a.Position.X),
		float32(a.Position.Y),
		float32(a.Size),
		color.White,
		false,
	)
}
