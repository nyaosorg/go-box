package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	box "github.com/zetamatta/go-box"
	"github.com/mattn/go-colorable"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
		return
	}
	list := strings.Split(string(data), "\n")
	for i := 0; i < len(list); i++ {
		list[i] = strings.TrimSpace(list[i])
	}
	result := box.Choice(
		list,
		colorable.NewColorableStderr())

	fmt.Println(result)
	os.Exit(0)
}
