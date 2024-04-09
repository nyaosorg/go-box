package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/nyaosorg/go-box/v2"
)

var optionIndex = flag.Bool("index", false, "print index as result")

func mains(args []string) error {
	data, err := io.ReadAll(os.Stdin)
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
	indexes, err := box.ChooseMultiX(list, console)
	if err != nil {
		return err
	}
	fmt.Fprintln(console)

	if *optionIndex {
		for _, index := range indexes {
			fmt.Println(index)
		}
	} else {
		if indexes == nil {
			return errors.New("canceled")
		}
		for _, index := range indexes {
			fmt.Println(list[index])
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if err := mains(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
