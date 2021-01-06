package mandelbrot

import (
	"math/cmplx"
)

// Mandelbrot computes the number of iterations necessary to determine whether
// the supplied complex number c and the Mandelbrot orbit sequence tends to
// infinity or not.
func Mandelbrot(c complex128, maxIterations int) int {
	var n int
	var z complex128
	for cmplx.Abs(z) <= 2 && n < maxIterations {
		z = z*z + c
		n++
	}
	return n
}
