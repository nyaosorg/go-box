package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/nyaosorg/go-box/v2"
)

var optionIndex = flag.Bool("index", false, "print index as result")

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
	indexes := box.ChooseMulti(list, console)
	fmt.Fprintln(console)

	if *optionIndex {
		if indexes != nil {
			for _, index := range indexes {
				fmt.Println(index)
			}
		}
	} else {
		if indexes != nil {
			for _, index := range indexes {
				fmt.Println(list[index])
			}
		} else {
			fmt.Println("canceled")
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if err := main1(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
