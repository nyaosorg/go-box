package box

import (
	"context"
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
