package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/Mamdasn/im2dhistgo"

)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("Usage: ./im2dhistgo <input>")
		os.Exit(1)
	}
	imagename := os.Args[1]

	twodhist := im2dhistgo.Im2dhist_file(imagename)

	fmt.Println(twodhist)
	fmt.Println("Done.")
}
