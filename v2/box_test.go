package box

import (
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	Print(nil, []string{
		"aaaa", "bbbb", "cccc", "fjdaksljflkdajfkljsalkfjdlkf",
		"jfkldsjflkjdsalkfjlkdsajflkajds",
		"fsdfsdf"}, os.Stdout)
}
