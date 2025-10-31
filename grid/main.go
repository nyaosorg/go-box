package grid

import (
	"regexp"
)

var reduntantColorChangePattern = regexp.MustCompile("(\x1B[^m]+m).*?(\x1B[^m]+m)")

func cutReduntantColorChange(s string) string {
	for {
		m := reduntantColorChangePattern.FindStringSubmatchIndex(s)
		if m == nil || len(m) <= 0 {
			return s
		}
		// all = s[m[0]:m[1]]
		first := s[m[2]:m[3]]
		second := s[m[4]:m[5]]

		if first == second {
			s = s[:m[4]] + s[m[5]:]
		} else {
			return s[:m[4]] + cutReduntantColorChange(s[m[4]:])
		}
	}
}
