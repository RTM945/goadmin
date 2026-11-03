package utils

import "testing"

func TestTitle(t *testing.T) {
	for in, want := range map[string]string{
		"hello world": "Hello World",
		"HELLO_WORLD": "Hello_World",
		"it's ok":     "It's Ok",
		"中文 abc":      "中文 Abc",
		"":            "",
	} {
		if got := Title(in); got != want {
			t.Errorf("Title(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestUuid(t *testing.T) {
	seen := map[string]bool{}
	for range 100 {
		s := Uuid(32)
		if len(s) != 32 || seen[s] {
			t.Fatalf("Uuid(32) = %q", s)
		}
		seen[s] = true
	}
}
