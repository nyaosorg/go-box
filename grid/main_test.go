package grid

import (
	"testing"
)

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
