package p

import "testing"

func BenchmarkIsCookieNameValid(b *testing.B) {
	jub0bs := func(b *testing.B) {
		for range b.N {
			for _, name := range names {
				t = IsCookieNameValid(name)
			}
		}
	}
	std := func(b *testing.B) {
		for range b.N {
			for _, name := range names {
				t = IsCookieNameValidStd(name)
			}
		}
	}
	b.ResetTimer()
	b.Run("v=jub0bs", jub0bs)
	b.Run("v=std", std)
}
