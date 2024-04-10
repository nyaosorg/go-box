package box

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestPrint(t *testing.T) {
	var buffer strings.Builder

	Print(context.TODO(), []string{
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

func TestCutReduntantColorChange(t *testing.T) {
	source := "\x1B[32;1m....\x1B[32;1m hogehoge"
	expect := "\x1B[32;1m.... hogehoge"

	actual := cutReduntantColorChange(source)
	if expect != actual {
		t.Fatalf("expect `%s` but `%s`", expect, actual)
	}

	source = "\x1B[32;1m....\x1B[33;1m hogehoge"
	expect = source // not change

	actual = cutReduntantColorChange(source)
	if expect != actual {
		t.Fatalf("expect `%s` but `%s`", expect, actual)
	}

	source = "\x1B[32;1m....\x1B[32;1m....\x1B[32;1m hogehoge"
	expect = "\x1B[32;1m........ hogehoge"

	actual = cutReduntantColorChange(source)
	if expect != actual {
		t.Fatalf("expect `%s` but `%s`", expect, actual)
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
		width:  80,
		height: 25,
		tty:    &TstTty{history: []string{"l", "l", "\n"}},
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
