package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/zetamatta/go-box"
)

func main1(args []string) error {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	data = box.AnsiCutter.ReplaceAll(data, []byte{})
	list := strings.Split(string(data), "\n")
	switch len(list) {
	case 0:
		return nil
	case 1:
		fmt.Println(strings.TrimSpace(list[0]))
		return nil
	}
	for i := 0; i < len(list); i++ {
		list[i] = strings.TrimSpace(list[i])
	}
	console := colorable.NewColorableStderr()
	result := box.Choice(list, console)
	fmt.Fprintln(console)

	fmt.Println(result)
	return nil
}

func main() {
	if err := main1(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
