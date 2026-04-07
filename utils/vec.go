package utils

import "math"

type Vec struct {
	X, Y float64
}

func (v Vec) Add(other Vec) Vec {
	return Vec{v.X + other.X, v.Y + other.Y}
}

func (v Vec) Sub(other Vec) Vec {
	return Vec{v.X - other.X, v.Y - other.Y}
}

func (v Vec) Scale(f float64) Vec {
	return Vec{v.X * f, v.Y * f}
}

func (v Vec) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vec) Dot(other Vec) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vec) Rotate(angle float64) Vec {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return Vec{X: v.X*cos - v.Y*sin, Y: v.X*sin + v.Y*cos}
}

func (v Vec) Normalized() Vec {
	l := v.Length()
	if l == 0 {
		return Vec{}
	}
	return v.Scale(1 / l)
}
