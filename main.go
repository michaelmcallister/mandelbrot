package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/michaelmcallister/mandelbrot/mandelbrot"
)

const width, height = 600, 400

const maxIterations = 100

func getColor(m int) color.RGBA {
	c0 := 255 - uint8(m*255/maxIterations)
	return color.RGBA{c0, c0, c0, 0xFF}
}

type MandelbrotViewer struct {
	rendering bool
	zoom      float64
	centerX   float64
	centerY   float64
}

func (v *MandelbrotViewer) xToReal(x int) float64 {
	minX := v.centerX - v.zoom/2
	centerX := minX + float64(x)/float64(width)*v.zoom
	return centerX
}

func (v *MandelbrotViewer) yToImaginary(y int) float64 {
	minY := v.centerY - v.zoom/2
	centerY := minY + float64(y)/float64(width)*v.zoom
	return centerY
}

func (v *MandelbrotViewer) zoomIn() {
	v.zoom -= 0.3
}

func (v *MandelbrotViewer) zoomOut() {
	v.zoom += 0.3
}

func (v *MandelbrotViewer) reset() {
	v.zoom = 3.0
	v.centerX = -0.72
	v.centerY = 0.45
}

func (v *MandelbrotViewer) pan(x, y float64) {
	v.centerX += x
	v.centerY += y
}

// Update don't do nuffink.
func (v *MandelbrotViewer) Update(screen *ebiten.Image) error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		v.zoomIn()
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		v.zoomOut()
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		v.reset()
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		v.pan(-0.05, 0.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		v.pan(0.05, 0.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		v.pan(0.0, -0.05)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		v.pan(0.0, 0.05)
	}
	return nil
}

// Draw displays the mandelbrot set.
func (v *MandelbrotViewer) Draw(screen *ebiten.Image) {
	v.rendering = true
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := complex(v.xToReal(x), v.yToImaginary(y))
			m := mandelbrot.Mandelbrot(c, maxIterations)
			screen.Set(x, y, getColor(m))
		}
	}
	v.rendering = false
}

// Layout takes the outside size (e.v., the window size) and returns the (logical) screen size.
func (v *MandelbrotViewer) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func main() {
	ebiten.SetWindowTitle("Mandelbrot")
	v := &MandelbrotViewer{
		zoom:    3.0,
		centerX: -0.72,
		centerY: 0.45,
	}
	if err := ebiten.RunGame(v); err != nil {
		panic(err)
	}
}
