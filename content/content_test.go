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