package utils

import "testing"

func TestArrayValue(t *testing.T) {
	type M = map[string]any

	root := M{
		"a":        1,
		"b":        M{"c": "x", "d": M{"e": 5}},
		"dot.path": "hit",
	}

	cases := []struct {
		name   string
		key    string
		expect any
	}{
		{"direct", "a", 1},
		{"nested", "b.c", "x"},
		{"deep", "b.d.e", 5},
		{"non-exist", "b.d.z", nil},
		{"nil-root", "x", nil},
		{"dot-key", "dot.path", "hit"},
	}

	for _, tc := range cases {
		got := ArrayValue(root, tc.key)
		if got != tc.expect {
			t.Fatalf("%s: expect %v, got %v", tc.name, tc.expect, got)
		}
	}

	// default value fallback
	if got := ArrayValue(root, "b.d.z", "def"); got != "def" {
		t.Fatalf("default fallback: expect def, got %v", got)
	}
	if got := ArrayValue(nil, "any", "def"); got != "def" {
		t.Fatalf("nil root fallback: expect def, got %v", got)
	}
	// non-map in path returns default
	if got := ArrayValue(M{"a": 1}, "a.x", "def"); got != "def" {
		t.Fatalf("non-map path: expect def, got %v", got)
	}
	// empty key returns default
	if got := ArrayValue(root, "", "def"); got != "def" {
		t.Fatalf("empty key: expect def, got %v", got)
	}
}
