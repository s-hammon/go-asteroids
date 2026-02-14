package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Shot struct {
	Id       int
	Position Vector2
	Velocity Vector2
	Lifespan int
	Size     float64
}

func NewShot(id int, pos Vector2, vel Vector2) *Shot {
	return &Shot{
		Id:       id,
		Position: pos,
		Velocity: vel,
		Lifespan: ShotLifespan,
		Size:     ShotSize,
	}
}

func (s *Shot) Update() {
	s.Position = s.Position.Add(s.Velocity)
	s.Lifespan--
}

func (s *Shot) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(s.Position.X-s.Size, s.Position.Y-s.Size)
	screen.DrawImage(shotImage, op)
}

func (s *Shot) Radius() float64 {
	return s.Size
}

func (s *Shot) Pos() Vector2 {
	return s.Position
}
