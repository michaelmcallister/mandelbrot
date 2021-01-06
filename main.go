package main

import (
	"fmt"
	"image/color"
	"image/color/palette"
	"os"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/michaelmcallister/mandelbrot/mandelbrot"
)

var waitGroup sync.WaitGroup

// default parameters for Mandelbrot.
const (
	zoomFactor        = 1.1
	iterationStep     = 5
	width             = 600
	height            = 400
	rMin              = -2.0
	rMax              = 1.0
	iMin              = -1.0
	iMax              = 1.8
	zoom              = width / (rMax - rMin)
	defaultIterations = 256
)

type mandelbrotViewer struct {
	maxIterations int
	rMax          float64
	rMin          float64
	iMin          float64
	iMax          float64
	zoom          float64
	displayDebug  bool
}

// interpolate is responsible for determining the new co-ordinates for the
// complex plane.
// See: https://stackoverflow.com/questions/41796832/smooth-zoom-with-mouse-in-mandelbrot-set-c
func interpolate(start, end, interpolation float64) float64 {
	return start + ((end - start) * interpolation)
}

// mouseLocation returns the location on the complex plane where the mouse
// pointer is currently at.
func (v *mandelbrotViewer) mouseLocation() (float64, float64) {
	mX, mY := ebiten.CursorPosition()
	mouseRe := float64(mX)/(width/(v.rMax-v.rMin)) + v.rMin
	mouseIm := float64(mY)/(height/(v.iMax-v.iMin)) + v.iMin

	return mouseRe, mouseIm
}

// color returns a colour based on the current value of m.
func (v *mandelbrotViewer) color(m int) color.Color {
	if m > len(palette.Plan9)-1 {
		return color.Black
	}
	if m <= 0 {
		return color.White
	}
	return palette.Plan9[m]
}

func (v *mandelbrotViewer) debugPrint(screen *ebiten.Image) {
	if !v.displayDebug {
		return
	}
	var sb strings.Builder
	x, y := v.mouseLocation()
	sb.WriteString(fmt.Sprintf("Location: %f, %f\n", x, y))
	sb.WriteString(fmt.Sprintf("Zoom: %f\n", v.zoom))
	sb.WriteString(fmt.Sprintf("Max Iterations: %d\n", v.maxIterations))

	// always return nil.
	ebitenutil.DebugPrint(screen, sb.String())
}

func (v *mandelbrotViewer) zoomIn() {
	v.zoom *= zoomFactor
}

func (v *mandelbrotViewer) zoomOut() {
	v.zoom /= zoomFactor
}

// reset sets the mandelbrot to how it was at the start.
func (v *mandelbrotViewer) reset() {
	v.maxIterations = defaultIterations
	v.rMin = rMin
	v.rMax = rMax
	v.iMin = iMin
	v.iMax = iMax
	v.zoom = zoom
}

func (v *mandelbrotViewer) increaseMaxIterations() {
	v.maxIterations += iterationStep
}

func (v *mandelbrotViewer) decreaseMaxIterations() {
	if v.maxIterations <= iterationStep {
		return
	}
	v.maxIterations -= iterationStep
}

// Update handles input that manipulates the complex plan.
func (v *mandelbrotViewer) Update(screen *ebiten.Image) error {
	// Click to zoom and pan.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		interpolation := 1.0 / zoomFactor
		mouseRe, mouseIm := v.mouseLocation()
		v.rMin = interpolate(mouseRe, v.rMin, interpolation)
		v.iMin = interpolate(mouseIm, v.iMin, interpolation)
		v.rMax = interpolate(mouseRe, v.rMax, interpolation)
		v.iMax = interpolate(mouseIm, v.iMax, interpolation)
		v.zoomIn()
	}

	// Increase/Decrease the max iterations.
	if ebiten.IsKeyPressed(ebiten.KeyEqual) {
		v.increaseMaxIterations()
	}
	if ebiten.IsKeyPressed(ebiten.KeyMinus) {
		v.decreaseMaxIterations()
	}

	// Zoom in and out based on the scroll wheel, preserving location on the
	// complex plane.
	_, dY := ebiten.Wheel()
	if dY < 0.0 {
		v.zoomOut()
	}
	if dY > 0.0 {
		v.zoomIn()
	}

	// Zoom out based on the scroll wheel, preserving location on the complex
	// plane.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		v.zoomOut()
	}

	// Reset back to factory defaults.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		v.reset()
	}

	// Exit the proram if 'q' is pressed.
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		os.Exit(0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		v.displayDebug = !v.displayDebug
	}
	return nil
}

// Draw displays the mandelbrot set.
func (v *mandelbrotViewer) Draw(screen *ebiten.Image) {
	pix := make([]byte, width*height*4)
	l := width * height
	for i := 0; i < l; i++ {
		waitGroup.Add(1)
		go func(i int) {
			defer waitGroup.Done()
			x := i % width
			y := i / width

			cx := float64(x)/v.zoom + v.rMin
			cy := float64(y)/v.zoom + v.iMin
			c := complex(cx, cy)
			m := mandelbrot.Mandelbrot(c, v.maxIterations)

			r, g, b, a := v.color(m).RGBA()
			pix[4*i] = byte(r)
			pix[4*i+1] = byte(g)
			pix[4*i+2] = byte(b)
			pix[4*i+3] = byte(a)
		}(i)
	}
	waitGroup.Wait()
	screen.ReplacePixels(pix)
	v.debugPrint(screen)
}

// Layout takes the outside size (e.v., the window size) and returns the (logical) screen size.
func (v *mandelbrotViewer) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func main() {
	ebiten.SetWindowTitle("Mandelbrot")
	ebiten.SetWindowSize(width, height)
	ebiten.SetMaxTPS(ebiten.UncappedTPS)

	v := &mandelbrotViewer{
		maxIterations: defaultIterations,
		rMin:          rMin,
		rMax:          rMax,
		iMin:          iMin,
		iMax:          iMax,
		zoom:          zoom,
		displayDebug:  true,
	}
	if err := ebiten.RunGame(v); err != nil {
		panic(err)
	}
}
