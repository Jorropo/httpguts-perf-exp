package p

import (
	"slices"
	"strings"
	"testing"
	"unicode/utf8"

	"golang.org/x/net/http/httpguts"
)

func TestIsTokenRune(t *testing.T) {
	var r rune
	for ; r <= utf8.MaxRune; r++ {
		want := httpguts.IsTokenRune(r)
		got := IsTokenRune(r)
		if got != want {
			t.Fatalf("IsTokenRune(%q): got %t; want %t", r, got, want)
		}
	}
}

func TestValidHeaderFieldName(t *testing.T) {
	for _, name := range names {
		want := httpguts.ValidHeaderFieldName(name)
		got := ValidHeaderFieldName(name)
		if got != want {
			t.Fatalf("ValidHeaderFieldName(%q): got %t; want %t", name, got, want)
		}
	}
}

var names = []string{
	``,
	`Accept-Charset`,
	`Accept-Encoding`,
	`Access-Control-Request-Headers`,
	`Access-Control-Request-Method`,
	`Connection`,
	`Content-Length`,
	`Cookie`,
	`Cookie2`,
	`Date`,
	`DNT`,
	`Expect`,
	`Host`,
	`Keep-Alive`,
	`Origin`,
	`Referer`,
	`Set-Cookie`,
	`TE`,
	`Trailer`,
	`Transfer-Encoding`,
	`Upgrade`,
	`Via`,
}

func init() { // augments names with their lowercase counterparts
	names = slices.Grow(names, len(names))
	for _, name := range names {
		names = append(names, strings.ToLower(name))
	}
}

var t bool

func BenchmarkValidHeaderFieldName(b *testing.B) {
	jub0bs := func(b *testing.B) {
		for range b.N {
			for _, name := range names {
				t = ValidHeaderFieldName(name)
			}
		}
	}
	std := func(b *testing.B) {
		for range b.N {
			for _, name := range names {
				t = httpguts.ValidHeaderFieldName(name)
			}
		}
	}
	b.ResetTimer()
	b.Run("v=jub0bs", jub0bs)
	b.Run("v=std", std)
}

func BenchmarkIsTokenRune(b *testing.B) {
	jub0bs := func(b *testing.B) {
		for range b.N {
			for _, name := range names {
				for _, r := range name {
					t = IsTokenRune(r)
				}
			}
		}
	}
	std := func(b *testing.B) {
		for range b.N {
			for _, name := range names {
				for _, r := range name {
					t = IsTokenRune(r)
				}
			}
		}
	}
	b.ResetTimer()
	b.Run("v=jub0bs", jub0bs)
	b.Run("v=std", std)
}
