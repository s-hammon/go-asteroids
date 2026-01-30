package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Asteroid struct {
	Id       int
	Position Vector2
	Velocity Vector2
	Rotation float64
	Size     float64
}

func NewAsteroid(id int, pos Vector2, vel Vector2, size float64) *Asteroid {
	return &Asteroid{
		Id:       id,
		Position: pos,
		Velocity: vel,
		Rotation: 0,
		Size:     size,
	}
}

func (a *Asteroid) Update() {
	a.Position = a.Position.Add(a.Velocity)
	a.Rotation += AsteroidRotationSpeed
}

func (a *Asteroid) Draw(screen *ebiten.Image) {
	img := GetCircleImage(a.Size)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(a.Position.X-a.Size, a.Position.Y-a.Size)
	screen.DrawImage(img, op)
}

func (a *Asteroid) Radius() float64 {
	return a.Size
}

func (a *Asteroid) Pos() Vector2 {
	return a.Position
}
