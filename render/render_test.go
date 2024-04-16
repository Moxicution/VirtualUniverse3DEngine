// Copyright © 2024 Galvanized Logic Inc.

package render

import (
	"testing"

	"github.com/gazed/vu/math/lin"
)

// ============================================================================
// Benchmarking

// go test -bench=64
// Check the cost of converting 4x4 matrix of float64 to float32
// The last few runs showed around 3ns
func Benchmark64to32(b *testing.B) {
	m32 := &m4{}
	m64 := &lin.M4{Xx: 11, Xy: 12, Xz: 13, Xw: 14, Yx: 21, Yy: 22, Yz: 23, Yw: 24, Zx: 31, Zy: 32, Zz: 33, Zw: 34, Wx: 41, Wy: 42, Wz: 43, Ww: 44}
	for cnt := 0; cnt < b.N; cnt++ {
		m32.set64(m64)
	}
}
