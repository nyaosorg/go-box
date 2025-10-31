package box

import (
	"context"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/nyaosorg/go-box/v3/internal/lazy"
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

var wtRuneWidth = lazy.Of[*runewidth.Condition]{
	New: func() *runewidth.Condition {
		c := runewidth.NewCondition()
		if os.Getenv("WT_SESSION") != "" && os.Getenv("WT_PROFILE_ID") != "" {
			c.EastAsianWidth = false
		}
		return c
	},
}

// Println outputs the given items in a tabular layout and appends
// a newline after the final line.
func Println(nodes []string, out io.Writer) error {
	b, err := New()
	if err != nil {
		return err
	}
	b.height = 0
	_, _, err = b.Println(nodes, 0, out)
	b.Close()
	return err
}

func truncate(s string, w int) string {
	return wtRuneWidth.Value().Truncate(strings.TrimSpace(s), w, "")
}

type nodeT struct {
	Index int
	Text  string
}

// SelectIndex returns the indexes that user selected.
func SelectIndex(sources []string, multi bool, out io.Writer) ([]int, error) {
	return SelectIndexContext(context.TODO(), sources, multi, out)
}

func SelectIndexContext(ctx context.Context, sources []string, multi bool, out io.Writer) ([]int, error) {
	b, err := New()
	if err != nil {
		return nil, err
	}
	r, err := b.SelectIndexContext(ctx, sources, multi, out)
	b.Close()
	return r, err
}

// SelectString returns the strings that user selected.
func SelectString(sources []string, multi bool, out io.Writer) ([]string, error) {
	return SelectStringContext(context.TODO(), sources, multi, out)
}

func SelectStringContext(ctx context.Context, sources []string, multi bool, out io.Writer) ([]string, error) {
	b, err := New()
	if err != nil {
		return nil, err
	}
	r, err := b.SelectStringContext(ctx, sources, multi, out)
	b.Close()
	return r, err
}
