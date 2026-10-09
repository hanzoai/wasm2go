package gcasm

import "testing"

// A named slot or argument whose package path holds a hyphen is matched whole,
// its offset kept: `github.com/hanzo-inc/...` once left `github.com/hanzo-` in
// front of every rewritten slot, which the assembler cannot read.
func TestA64SlotsInHyphenatedPaths(t *testing.T) {
	for _, c := range []struct {
		re   string
		text string
		off  string
	}{
		{"sp", "github.com/hanzo-inc/zen/host/core/p0.autotmp_12-128(SP)", "-128"},
		{"sp", "github.com/hanzo-inc/zen/host/core/p0.~r0+8(SP)", "+8"},
		{"sp", "pkg.x-16(SP)", "-16"},
		{"sp", "github.com/a-b/c.d(SP)", ""},
		{"fp", "github.com/hanzo-inc/zen/host/core/p0.m+0(FP)", "+0"},
		{"fp", "github.com/hanzo-inc/zen/host/core/p0.l0+8(FP)", "+8"},
	} {
		re := a64NamedSlotRe
		if c.re == "fp" {
			re = a64FPRe
		}
		m := re.FindStringSubmatch("\tMOVD\tR0, " + c.text)
		if m == nil || m[0] != c.text || m[1] != c.off {
			t.Errorf("%s %q: matched %q", c.re, c.text, m)
		}
	}
}
