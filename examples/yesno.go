//go:build run

package main

import (
	"os"

	"github.com/nyaosorg/go-box/v2"
)

func main() {
	println("Are you sure ?")
	choose, err := box.SelectString([]string{"Yes", "No"}, false, os.Stderr)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	println()
	if len(choose) >= 1 {
		println("You selected ->", choose[0])
	} else {
		println("You did not select any items")
	}
}
