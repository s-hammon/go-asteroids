package main

import "math"

func FromAngleMagnitude(angle, magnitude float64) Vector2 {
	sin, cos := math.Sincos(angle)
	return Vector2{
		X: cos * magnitude,
		Y: sin * magnitude,
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

func (v Vector2) DistanceTo(other Vector2) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (v Vector2) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

func (v Vector2) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func isCircleOffscreen(pos Vector2, size float64) bool {
	return pos.X+size < 0 ||
		pos.X-size > ScreenWidth ||
		pos.Y+size < 0 ||
		pos.Y-size > ScreenHeight
}
