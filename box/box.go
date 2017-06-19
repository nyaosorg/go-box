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
	switch len(list) {
	case 0:
		os.Exit(1)
	case 1:
		fmt.Println(strings.TrimSpace(list[0]))
		os.Exit(0)
	}
	for i := 0; i < len(list); i++ {
		list[i] = strings.TrimSpace(list[i])
	}
	result := box.Choice(
		list,
		colorable.NewColorableStderr())

	fmt.Println(result)
	os.Exit(0)
}
