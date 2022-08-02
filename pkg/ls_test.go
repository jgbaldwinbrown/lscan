package lscan

import (
	"testing"
	"strings"
	"golang.org/x/exp/slices"
)

func cmpStrSlSl(sls1, sls2 [][]string) bool {
	if len(sls1) != len(sls2) {
		return false
	}
	for i, sl1 := range sls1 {
		if len(sl1) != len(sls2[i]) {
			return false
		}
		for j, s1 := range sl1 {
			if s1 != sls2[i][j] {
				return false
			}
		}
	}
	return true
}

func TestScanner(t *testing.T) {
	expected := [][]string {
		[]string { "apple", "banana" },
		[]string { "1", "2", "3" },
	}
	in := "apple\tbanana\n1\t2\t3\n"

	s := NewScanner(strings.NewReader(in), ByByte('\t'))
	var out [][]string
	for s.Scan() {
		out = append(out, slices.Clone(s.Line()))
	}

	if !cmpStrSlSl(expected, out) {
		t.Errorf("expected does not equal out: %v; %v", expected, out)
	}
}
