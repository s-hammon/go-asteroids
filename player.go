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

type PlayerState int

const (
	PlayerStateIdle PlayerState = iota
	PlayerStateRespawn
	PlayerStateDead
)

type Player struct {
	Position Vector2
	Velocity Vector2
	Rotation float64
	Size     float64
	State    PlayerState

	shootCooldown   float64
	respawnCooldown float64
	lives           int
	blinkTimer      int
}

func NewPlayer() *Player {
	return &Player{
		Position: Vector2{
			X: ScreenWidthH,
			Y: ScreenHeightH,
		},
		Rotation:        -math.Pi / 2,
		Size:            PlayerSize,
		State:           PlayerStateIdle,
		shootCooldown:   0,
		respawnCooldown: 0,
		lives:           PlayerInitialLives,
	}
}

func (p *Player) Update() {
	p.shootCooldown = max(0, p.shootCooldown-1)
	p.respawnCooldown = max(0, p.respawnCooldown-1)

	if p.respawnCooldown == 0 {
		p.blinkTimer = 0
		p.State = PlayerStateIdle
	} else {
		p.blinkTimer++
		if p.blinkTimer >= RespawnBlinkRate {
			p.blinkTimer = 0
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.Rotation -= PlayerRotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.Rotation += PlayerRotationSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		sin, cos := math.Sincos(p.Rotation)
		thrust := Vector2{
			X: cos,
			Y: sin,
		}.Scale(PlayerSpeed)
		p.Velocity = p.Velocity.Add(thrust)
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) && p.shootCooldown == 0 {
		p.Shoot(shootCh)
		p.shootCooldown = PlayerShootCooldown
	}

	p.Velocity = p.Velocity.Scale(PlayerFriction)
	p.Position = p.Position.Add(p.Velocity)
	p.Position = p.Position.Clamp(0, 0, ScreenWidth, ScreenHeight)
}

func (p *Player) Draw(screen *ebiten.Image) {
	if p.respawnCooldown > 0 && p.blinkTimer < RespawnBlinkRate/2 {
		return
	}

	img := GetPlayerImage()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-PlayerSize*2, -PlayerSize*2)
	op.GeoM.Rotate(p.Rotation + math.Pi/2)
	op.GeoM.Translate(p.Position.X, p.Position.Y)
	screen.DrawImage(img, op)
}

func (p *Player) Radius() float64 {
	return PlayerSize
}

func (p *Player) Pos() Vector2 {
	return p.Position
}

func (p *Player) Shoot(ch chan<- struct{}) {
	ch <- struct{}{}
}

func (p *Player) respawn() {
	p.State = PlayerStateRespawn
	p.respawnCooldown = PlayerRespawnCooldown
}
