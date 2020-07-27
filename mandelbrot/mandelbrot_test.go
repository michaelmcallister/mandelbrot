package mandelbrot

import "testing"

func TestMandelbrot(t *testing.T) {
	const maxIterations = 80
	testCases := []struct {
		desc  string
		input complex128
		want  int
	}{
		{
			desc:  "-1,-1i = 3",
			input: complex(-1.0, -1.0),
			want:  3,
		},
		{
			desc:  "-1,-0.5i = 5",
			input: complex(-1.0, -0.5),
			want:  5,
		},
		{
			desc:  "-1, +0i = 80",
			input: complex(-1.0, 0.0),
			want:  80,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := Mandelbrot(tC.input, maxIterations)
			if got != tC.want {
				t.Errorf("Mandelbrot(%v, %d) = %d, want=%d", tC.input, maxIterations, got, tC.want)
			}
		})
	}
}

func benchmarkMandelbrot(i int, b *testing.B) {
	c := complex(-1.0, 0.0)
	for n := 0; n < b.N; n++ {
		Mandelbrot(c, i)
	}
}

func BenchmarkMandelbrot10(b *testing.B)   { benchmarkMandelbrot(10, b) }
func BenchmarkMandelbrot50(b *testing.B)   { benchmarkMandelbrot(50, b) }
func BenchmarkMandelbrot100(b *testing.B)  { benchmarkMandelbrot(100, b) }
func BenchmarkMandelbrot1000(b *testing.B) { benchmarkMandelbrot(1000, b) }
