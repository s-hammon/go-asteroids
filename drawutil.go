package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	asteroidImageCache     = make(map[float64]*ebiten.Image)
	playerImage, shotImage *ebiten.Image
)

func init() {
	shotImage = initShotImage()
}

func initShotImage() *ebiten.Image {
	diameter := ShotSize * 2
	img := ebiten.NewImage(int(diameter+2), int(diameter+2))

	vector.FillCircle(
		img,
		float32(ShotSize+1),
		float32(ShotSize+1),
		float32(ShotSize),
		dragonRed,
		false,
	)

	return img

}

func GetAsteroidImage(size float64) *ebiten.Image {
	img, ok := asteroidImageCache[size]
	if ok {
		return img
	}

	diameter := size * 2
	img = ebiten.NewImage(int(diameter+2), int(diameter+2))

	vector.FillCircle(
		img,
		float32(size+1),
		float32(size+1),
		float32(size),
		winterYellow,
		false,
	)

	asteroidImageCache[size] = img
	return img
}

func GetPlayerImage() *ebiten.Image {
	if playerImage != nil {
		return playerImage
	}

	size := PlayerSize * 4
	img := ebiten.NewImage(size, size)

	centerX := float64(size / 2)
	centerY := float64(size / 2)

	sin, cos := math.Sincos(-math.Pi / 2)
	j := Vector2{
		X: centerX + cos*PlayerSize*2,
		Y: centerY + sin*PlayerSize*2,
	}

	sin, cos = math.Sincos(-math.Pi/2 + 2*math.Pi/3)
	k := Vector2{
		X: centerX + cos*PlayerSize,
		Y: centerY + sin*PlayerSize,
	}

	sin, cos = math.Sincos(-math.Pi/2 - 2*math.Pi/3)
	l := Vector2{
		X: centerX + cos*PlayerSize,
		Y: centerY + sin*PlayerSize,
	}

	r := float32(dragonGreen.R) / 255.0
	g := float32(dragonGreen.G) / 255.0
	b := float32(dragonGreen.B) / 255.0
	a := float32(dragonGreen.A) / 255.0
	vertices := []ebiten.Vertex{
		{DstX: float32(j.X), DstY: float32(j.Y), ColorR: r, ColorG: g, ColorB: b, ColorA: a},
		{DstX: float32(k.X), DstY: float32(k.Y), ColorR: r, ColorG: g, ColorB: b, ColorA: a},
		{DstX: float32(l.X), DstY: float32(l.Y), ColorR: r, ColorG: g, ColorB: b, ColorA: a},
	}

	indices := []uint16{0, 1, 2}
	img.DrawTriangles(vertices, indices, whitePx, nil)

	playerImage = img
	return playerImage
}
