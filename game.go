package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand/v2"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	f *text.GoTextFace

	livesOp, scoreOp *text.DrawOptions

	shootCh chan struct{}
)

func init() {
	faceSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	f = &text.GoTextFace{
		Source: faceSource,
		Size:   24,
	}

	livesOp = &text.DrawOptions{}
	livesOp.ColorScale.ScaleWithColor(dragonRed)
	livesOp.GeoM.Translate(10, ScreenHeight-60)

	scoreOp = &text.DrawOptions{}
	scoreOp.ColorScale.ScaleWithColor(dragonRed)
	scoreOp.GeoM.Translate(10, ScreenHeight-30)

	shootCh = make(chan struct{}, 1)
}

type GameState int

const (
	GameStateTitle GameState = iota
	GameStatePlaying
	GameStatePaused
	GameStateOver
)

type Game struct {
	Player             *Player
	Asteroids          map[int]*Asteroid
	Shots              map[int]*Shot
	AsteroidSpawnTimer int
	State              GameState

	nextAsteroidId int32
	nextShotId     int32
	points         points
	pauseDelay     float64
}

func NewGame() *Game {
	return &Game{
		Player:             NewPlayer(),
		Asteroids:          make(map[int]*Asteroid),
		Shots:              make(map[int]*Shot),
		AsteroidSpawnTimer: AsteroidSpawnInterval,
		State:              GameStateTitle,
	}
}

func (g *Game) Update() error {
	g.pauseDelay = max(0, g.pauseDelay-1)

	switch g.State {
	default:
		panic(fmt.Sprintln("unknown game state:", g.State))
	case GameStateOver:
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.Restart()
		}
		if ebiten.IsKeyPressed(ebiten.KeyQ) {
			return ebiten.Termination
		}
	case GameStateTitle:
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.State = GameStatePlaying
		}
	case GameStatePaused:
		if ebiten.IsKeyPressed(ebiten.KeyP) {
			g.Resume()
		}
		if ebiten.IsKeyPressed(ebiten.KeyQ) {
			g.NewGame()
		}
	case GameStatePlaying:
		if ebiten.IsKeyPressed(ebiten.KeyP) {
			g.Pause()
			return nil
		}

		g.Player.Update()

		select {
		default:
		case <-shootCh:
			g.spawnNewShot(g.Player)
		}

		g.AsteroidSpawnTimer--
		if g.AsteroidSpawnTimer <= 0 {
			if len(g.Asteroids) <= AsteroidMaxCount {
				g.spawnNewAsteroid()
			}
			g.AsteroidSpawnTimer = AsteroidSpawnInterval
		}

		for id, a := range g.Asteroids {
			a.Update()
			if Collides(a, g.Player) {
				if g.Player.State == PlayerStateRespawn {
					continue
				}
				g.Player.lives--
				if isOutOfLives(g.Player) {
					g.Player.State = PlayerStateDead
					g.State = GameStateOver
					return nil
				}

				g.Player.respawn()
			}

			if isCircleOffscreen(a.Position, a.Size) {
				delete(g.Asteroids, id)
			}
		}

		for id, s := range g.Shots {
			s.Update()
			for _, a := range g.Asteroids {
				if Collides(a, s) {
					point := g.splitAsteroid(a)
					g.updatePoints(point)
					delete(g.Shots, id)
					continue
				}
			}
			if isCircleOffscreen(s.Position, s.Size) {
				delete(g.Shots, id)
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	switch g.State {
	default:
		panic(fmt.Sprintln("unknown game state:", g.State))
	case GameStateTitle:
		op := &text.DrawOptions{}
		op.ColorScale.ScaleWithColor(springGreen)
		op.GeoM.Translate(ScreenWidthH-60, ScreenHeightH-90)
		text.Draw(screen, "ASTEROIDS", f, op)

		op.ColorScale.Reset()
		op.ColorScale.ScaleWithColor(color.White)
		op.GeoM.Translate(-30, 60)
		text.Draw(screen, "Press [S] to start...", f, op)
	case GameStateOver:
		op := &text.DrawOptions{}
		op.ColorScale.ScaleWithColor(dragonRed)
		op.GeoM.Translate(ScreenWidthH-60, ScreenHeightH-30)
		text.Draw(screen, "GAME OVER!", f, op)

		op.ColorScale.Reset()
		op.ColorScale.ScaleWithColor(color.White)
		op.GeoM.Translate(-30, 30)
		text.Draw(screen, "Press [R] to restart", f, op)

		op.GeoM.Translate(0, 30)
		text.Draw(screen, "Press [Q] to quit...", f, op)
	case GameStatePaused:
		op := &text.DrawOptions{}
		op.ColorScale.ScaleWithColor(color.White)
		op.GeoM.Translate(ScreenWidthH-100, ScreenHeightH-30)
		text.Draw(screen, "Press [P] to continue...", f, op)

		op.GeoM.Translate(0, 30)
		text.Draw(screen, "Press [Q] to quit...", f, op)
	case GameStatePlaying:
		for _, a := range g.Asteroids {
			a.Draw(screen)
		}
		for _, s := range g.Shots {
			s.Draw(screen)
		}

		g.Player.Draw(screen)
		text.Draw(screen, fmt.Sprintf("LIVES: %d", g.Player.lives), f, livesOp)
		text.Draw(screen, fmt.Sprintf("SCORE: %d", g.points.total), f, scoreOp)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) NewGame() {
	g.Player = NewPlayer()
	g.Asteroids = make(map[int]*Asteroid)
	g.Shots = make(map[int]*Shot)
	g.AsteroidSpawnTimer = AsteroidSpawnInterval
	g.State = GameStateTitle
	g.nextAsteroidId = 0
	g.nextShotId = 0
	g.points = points{}
	g.pauseDelay = 24
}

func (g *Game) Restart() {
	g.Player = NewPlayer()
	g.Asteroids = make(map[int]*Asteroid)
	g.Shots = make(map[int]*Shot)
	g.AsteroidSpawnTimer = AsteroidSpawnInterval
	g.State = GameStatePlaying
	g.nextAsteroidId = 0
	g.nextShotId = 0
	g.points = points{}
	g.pauseDelay = 24
}

func (g *Game) Pause() {
	if g.pauseDelay > 0 {
		return
	}

	g.pauseDelay = 24
	g.State = GameStatePaused
}

func (g *Game) Resume() {
	if g.pauseDelay > 0 {
		return
	}

	g.pauseDelay = 24
	g.State = GameStatePlaying
}

func (g *Game) spawnNewAsteroid() {
	var (
		pos         Vector2
		baseAngle   float64
		spreadRange = math.Pi / 3
	)
	switch side := rand.IntN(4); side {
	case 0:
		pos = Vector2{X: float64(rand.IntN(ScreenWidth)), Y: 0}
		baseAngle = math.Pi / 2
	case 1:
		pos = Vector2{X: ScreenWidth, Y: float64(rand.IntN(ScreenHeight))}
		baseAngle = math.Pi
	case 2:
		pos = Vector2{X: float64(rand.IntN(ScreenWidth)), Y: ScreenHeight}
		baseAngle = 3 * math.Pi / 2
	case 3:
		pos = Vector2{X: 0, Y: float64(rand.IntN(ScreenHeight))}
		baseAngle = 0
	}

	speed := AsteroidMinSpeed + rand.Float64()*(AsteroidMaxSpeed-AsteroidMinSpeed)
	randomOffset := rand.Float64()*spreadRange - (spreadRange / 2)
	angle := baseAngle + randomOffset

	sin, cos := math.Sincos(angle)
	vel := Vector2{X: cos * speed, Y: sin * speed}
	kind := 1 + rand.IntN(AsteroidKinds)

	id := getNextId(&g.nextAsteroidId)
	g.Asteroids[id] = NewAsteroid(id, pos, vel, AsteroidMinSize*float64(kind))
}

func (g *Game) spawnNewShot(p *Player) {
	sin, cos := math.Sincos(p.Rotation)
	pos := Vector2{
		X: p.Position.X + cos*(p.Size*2),
		Y: p.Position.Y + sin*(p.Size*2),
	}

	sin, cos = math.Sincos(p.Rotation)
	vel := Vector2{
		X: p.Velocity.X + cos*ShotSpeed,
		Y: p.Velocity.Y + sin*ShotSpeed,
	}

	id := getNextId(&g.nextShotId)
	g.Shots[id] = NewShot(id, pos, vel)
}

const (
	pointLarge int = iota + 1
	pointMedium
	pointSmall
)

func (g *Game) splitAsteroid(a *Asteroid) int {
	delete(g.Asteroids, a.Id)

	if a.Size <= AsteroidMinSize {
		return pointSmall
	}

	newSize := a.Size - AsteroidMinSize

	originalAngle := a.Velocity.Angle()
	originalMagnitude := a.Velocity.Magnitude()

	angleOffset := (20 + rand.Float64()*30) * (math.Pi / 180)

	angle1 := originalAngle + angleOffset
	vel1 := FromAngleMagnitude(angle1, originalMagnitude*1.2)
	id1 := getNextId(&g.nextAsteroidId)
	g.Asteroids[id1] = NewAsteroid(id1, a.Position, vel1, newSize)

	angle2 := originalAngle - angleOffset
	vel2 := FromAngleMagnitude(angle2, originalMagnitude*1.2)
	id2 := getNextId(&g.nextAsteroidId)
	g.Asteroids[id2] = NewAsteroid(id2, a.Position, vel2, newSize)

	if a.Size >= AsteroidMaxSize {
		return pointLarge
	} else {
		return pointMedium
	}
}

type points struct {
	total         int
	medium, large int
}

func (g *Game) updatePoints(point int) {
	switch point {
	case pointSmall:
		g.points.total++
	case pointMedium:
		if g.points.medium >= 1 {
			g.points.total++
			g.points.medium = 0
		} else {
			g.points.medium++
		}
	case pointLarge:
		if g.points.large >= 3 {
			g.points.total++
			g.points.medium = 0
		} else {
			g.points.large++
		}
	}
}

func getNextId(id *int32) int {
	return int(atomic.AddInt32(id, 1) - 1)
}

func isOutOfLives(p *Player) bool {
	return p.lives <= 0
}

func Collides(a, b CircleActor) bool {
	distance := a.Pos().DistanceTo(b.Pos())
	return distance <= a.Radius()+b.Radius()
}
