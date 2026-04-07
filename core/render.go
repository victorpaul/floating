package core

import (
	"image"
	"image/color"
	"math"

	"kosh/vpaul/floating/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// whiteSubImage is a 1x1 white pixel used as a texture source for path rendering.
var whiteSubImage *ebiten.Image

func getWhiteSubImage() *ebiten.Image {
	if whiteSubImage == nil {
		img := ebiten.NewImage(3, 3)
		img.Fill(color.White)
		whiteSubImage = img.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	}
	return whiteSubImage
}

// drawRotatedFilledRect draws a filled rectangle centered at (cx, cy) with given
// half-extents and rotation angle (radians).
func DrawRotatedFilledRect(screen *ebiten.Image, cx, cy float32, half utils.Vec, angle float64, clr color.RGBA) {
	vs, is := rotatedRectVertices(cx, cy, half, angle)
	applyColor(vs, clr)
	screen.DrawTriangles(vs, is, getWhiteSubImage(), &ebiten.DrawTrianglesOptions{AntiAlias: true})
}

// StrokeRotatedRect draws a stroked (outline) rectangle centered at (cx, cy).
func StrokeRotatedRect(screen *ebiten.Image, cx, cy float32, half utils.Vec, angle float64, strokeWidth float32, clr color.RGBA) {
	path := rotatedRectPath(cx, cy, half, angle)
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{Width: strokeWidth})
	applyColor(vs, clr)
	screen.DrawTriangles(vs, is, getWhiteSubImage(), &ebiten.DrawTrianglesOptions{AntiAlias: true})
}

func rotatedRectPath(cx, cy float32, half utils.Vec, angle float64) vector.Path {
	cos32 := float32(math.Cos(angle))
	sin32 := float32(math.Sin(angle))
	hw, hh := float32(half.X), float32(half.Y)

	corners := [4][2]float32{
		{-hw, -hh}, {hw, -hh}, {hw, hh}, {-hw, hh},
	}

	var path vector.Path
	for i, c := range corners {
		wx := cx + c[0]*cos32 - c[1]*sin32
		wy := cy + c[0]*sin32 + c[1]*cos32
		if i == 0 {
			path.MoveTo(wx, wy)
		} else {
			path.LineTo(wx, wy)
		}
	}
	path.Close()
	return path
}

func rotatedRectVertices(cx, cy float32, half utils.Vec, angle float64) ([]ebiten.Vertex, []uint16) {
	path := rotatedRectPath(cx, cy, half, angle)
	return path.AppendVerticesAndIndicesForFilling(nil, nil)
}

func applyColor(vs []ebiten.Vertex, clr color.RGBA) {
	cr := float32(clr.R) / 255
	cg := float32(clr.G) / 255
	cb := float32(clr.B) / 255
	ca := float32(clr.A) / 255
	for i := range vs {
		vs[i].SrcX, vs[i].SrcY = 1, 1
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = cr, cg, cb, ca
	}
}
