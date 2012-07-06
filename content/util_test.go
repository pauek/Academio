package content

import (
	"testing"
)

func TestWalk(t *testing.T) {
	roots = []string{"/u", "/o", "/v", "/r/a/b"}
	cases := []struct {
		abs, rel string
		nlev int
	}{
		{"/u/bin/a", "bin/a", 2},
		{"/r/a/b/c", "c", 1},
		{"/v", "", 0},
		{"/o/xyz/t", "xyz/t", 2},
		{"/r/a/b/c/d/e", "c/d/e", 3},
		{"/u/1/2/3/4", "1/2/3/4", 4},
	}
	for _, c := range cases {
		rel := removeRoot(c.abs)
		if rel != c.rel {
			t.Errorf("'%s' != '%s'", c.rel, rel)
		}
		nlev := numLevels(rel)
		if nlev != c.nlev {
			t.Errorf("%d != %d", c.nlev, nlev)
		}
	}
}

func TestRemoveOrder(t *testing.T) {
	cases := [][]string{
		{"a/b/c", "a/b/c"},
		{"1. a/2. b/3. c", "a/b/c"},
		{"001. a/002. b/003. c", "a/b/c"},
		{"123. a/0. c", "a/c"},
		{"0. xxxx", "xxxx"},
		{"10298347. añsldkfj/01. pqwoieru", "añsldkfj/pqwoieru"},
		{"1x a/2x b", "1x a/2x b"},
		{"a1. /b2. /c3. ", "a1. /b2. /c3. "},
		{"123/456/789", "123/456/789"},
	}
	for _, c := range cases {
		x := removeOrder(c[0])
		if x != c[1] {
			t.Errorf("%q != %q", x, c[1])
		}
	}
}

func TestParentID(t *testing.T) {
	cases := [][]string{
		{"a.b.c", "a.b"},
		{"a", ""},
		{"a.b", "a"},
		{"aaa.bbb.ccc", "aaa.bbb"},
		{"", ""},
	}
	for _, c := range cases {
		pid := parentID(c[0])
		if pid != c[1] {
			t.Errorf("%q != %q", pid, c[1])
		}
	}
}