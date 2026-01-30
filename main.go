package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth   = 1280
	ScreenWidthH  = ScreenWidth / 2
	ScreenHeight  = 720
	ScreenHeightH = ScreenHeight / 2

	PlayerRotationSpeed   = 0.05
	PlayerSpeed           = 0.15
	PlayerFriction        = 0.97
	PlayerSize            = 15
	PlayerShootCooldown   = 24
	PlayerInitialLives    = 3
	PlayerRespawnCooldown = 180
	RespawnBlinkRate      = 20

	AsteroidKinds         = 3
	AsteroidMinSize       = 20.0
	AsteroidMaxSize       = 60.0
	AsteroidMinSpeed      = 0.5
	AsteroidMaxSpeed      = 2.0
	AsteroidSpawnInterval = 35
	AsteroidRotationSpeed = 0.02
	AsteroidSides         = 12
	AsteroidMaxCount      = 30

	ShotLifespan = 2
	ShotSize     = 5
	ShotSpeed    = 10
)

const (
	statePlaying int = iota
	stateGameOver
)

func main() {
	ebiten.SetWindowTitle("Ye Olde Asteroids")
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
