package box

import (
	"io"
	"strings"
	"testing"

	"github.com/nyaosorg/go-box/v3/grid"
)

func TestPrint(t *testing.T) {
	var buffer strings.Builder

	Println([]string{
		"aaaa", "bbbb", "cccc", "fjdaksljflkdajfkljsalkfjdlkf",
		"jfkldsjflkjdsalkfjlkdsajflkajds",
		"fsdfsdf"}, &buffer)

	actual := buffer.String()
	expect := `aaaa                            fjdaksljflkdajfkljsalkfjdlkf
bbbb                            jfkldsjflkjdsalkfjlkdsajflkajds
cccc                            fsdfsdf
`
	if actual != expect {
		t.Fatalf("expect `%s` buf `%s`", expect, actual)
	}

}

type TstTty struct {
	history []string
}

func (t *TstTty) GetKey() (string, error) {
	if len(t.history) <= 0 {
		return "", io.EOF
	}
	result := t.history[0]
	t.history = t.history[1:]
	return result, nil
}

func (t *TstTty) Close() error {
	return nil
}

func TestSelectIndex(t *testing.T) {
	b := &Box{
		Grid: grid.Grid{
			Width:  80,
			Height: 25,
		},
		tty: &TstTty{history: []string{"l", "l", "\n"}},
	}
	list := []string{"A", "B", "C", "D", "E"}
	r, err := b.SelectIndex(list, false, io.Discard)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(r) != 1 {
		t.Fatalf("too few result: %d", len(r))
	}
	if r[0] != 2 {
		t.Fatalf("expect %v,but %v", 2, r[0])
	}
}
