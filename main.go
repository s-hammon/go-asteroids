package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 720

	PlayerRotationSpeed = 0.05
	PlayerSpeed         = 0.15
	PlayerFriction      = 0.96
	PlayerSize          = 15

	AsteroidMinSize       = 20.0
	AsteroidMaxSize       = 60.0
	AsteroidMinSpeed      = 0.5
	AsteroidMaxSpeed      = 2.0
	AsteroidSpawnInterval = 35
	AsteroidRotationSpeed = 0.02
	AsteroidSides         = 12
	AsteroidMaxCount      = 30
)

type Game struct {
	Player             *Player
	Asteroids          []*Asteroid
	AsteroidSpawnTimer int
}

func (g *Game) Update() error {
	g.Player.Update()

	g.AsteroidSpawnTimer--
	if g.AsteroidSpawnTimer <= 0 {
		if len(g.Asteroids) <= AsteroidMaxCount {
			g.spawnNewAsteroid()
		}
		g.AsteroidSpawnTimer = AsteroidSpawnInterval
	}

	for _, a := range g.Asteroids {
		a.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	g.Player.Draw(screen)

	for _, a := range g.Asteroids {
		a.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	if whitePx == nil {
		log.Fatalln("whitePx not set!")
	}
	ebiten.SetWindowTitle("Ye Olde Asteroids")
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	game := &Game{
		Player:             NewPlayer(),
		Asteroids:          make([]*Asteroid, 0),
		AsteroidSpawnTimer: AsteroidSpawnInterval,
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Vector2 struct {
	X, Y float64
}

func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vector2) Scale(factor float64) Vector2 {
	return Vector2{
		X: v.X * factor,
		Y: v.Y * factor,
	}
}

func (v Vector2) Clamp(x0, y0, x1, y1 float64) Vector2 {
	v.X = max(x0, min(v.X, x1))
	v.Y = max(y0, min(v.Y, y1))
	return v
}

func (g *Game) spawnNewAsteroid() {
	var (
		pos         Vector2
		baseAngle   float64
		spreadRange = math.Pi / 3
	)
	switch side := rand.Intn(4); side {
	case 0:
		pos = Vector2{X: float64(rand.Intn(ScreenWidth)), Y: 0}
		baseAngle = math.Pi / 2
	case 1:
		pos = Vector2{X: ScreenWidth, Y: float64(rand.Intn(ScreenHeight))}
		baseAngle = math.Pi
	case 2:
		pos = Vector2{X: float64(rand.Intn(ScreenWidth)), Y: ScreenHeight}
		baseAngle = 3 * math.Pi / 2
	case 3:
		pos = Vector2{X: 0, Y: float64(rand.Intn(ScreenHeight))}
		baseAngle = 0
	}

	speed := AsteroidMinSpeed + rand.Float64()*(AsteroidMaxSpeed-AsteroidMinSpeed)
	randomOffset := rand.Float64()*spreadRange - (spreadRange / 2)
	angle := baseAngle + randomOffset

	vel := Vector2{X: math.Cos(angle) * speed, Y: math.Sin(angle) * speed}
	size := AsteroidMinSize + rand.Float64()*(AsteroidMaxSize-AsteroidMinSize)
	g.Asteroids = append(g.Asteroids, NewAsteroid(pos, vel, size))
}
