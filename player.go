package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var whitePx *ebiten.Image

func init() {
	whitePx = ebiten.NewImage(1, 1)
	whitePx.Fill(color.White)
}

type Player struct {
	Position Vector2
	Velocity Vector2
	Rotation float64
}

func NewPlayer() *Player {
	return &Player{
		Position: Vector2{
			X: ScreenWidth / 2,
			Y: ScreenHeight / 2,
		},
		Rotation: -math.Pi / 2,
	}
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.Rotation -= PlayerRotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.Rotation += PlayerRotationSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		thrust := Vector2{
			X: math.Cos(p.Rotation),
			Y: math.Sin(p.Rotation),
		}.Scale(PlayerSpeed)
		p.Velocity = p.Velocity.Add(thrust)
	}

	p.Velocity = p.Velocity.Scale(PlayerFriction)
	p.Position = p.Position.Add(p.Velocity)
	p.Position = p.Position.Clamp(0, 0, ScreenWidth, ScreenHeight)
}

func (p *Player) Draw(screen *ebiten.Image) {
	a := Vector2{
		X: p.Position.X + math.Cos(p.Rotation)*PlayerSize*2,
		Y: p.Position.Y + math.Sin(p.Rotation)*PlayerSize*2,
	}

	b := Vector2{
		X: p.Position.X + math.Cos(p.Rotation+2*math.Pi/3)*PlayerSize,
		Y: p.Position.Y + math.Sin(p.Rotation+2*math.Pi/3)*PlayerSize,
	}

	c := Vector2{
		X: p.Position.X + math.Cos(p.Rotation-2*math.Pi/3)*PlayerSize,
		Y: p.Position.Y + math.Sin(p.Rotation-2*math.Pi/3)*PlayerSize,
	}

	vertices := []ebiten.Vertex{
		{DstX: float32(a.X), DstY: float32(a.Y), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(b.X), DstY: float32(b.Y), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(c.X), DstY: float32(c.Y), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
	indices := []uint16{0, 1, 2}

	screen.DrawTriangles(vertices, indices, whitePx, nil)
}
