package box_test

import (
	"strings"
	"testing"

	"github.com/zetamatta/go-box/v2"
)

func TestPrint(t *testing.T) {
	var buffer strings.Builder

	box.Print(nil, []string{
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
