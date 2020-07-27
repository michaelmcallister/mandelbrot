package mandelbrot

import "math/cmplx"

func Mandelbrot(c complex128, maxIterations int) int {
	var n int
	var z complex128
	for cmplx.Abs(z) <= 2 && n < maxIterations {
		z = z*z + c
		n++
	}
	return n
}
