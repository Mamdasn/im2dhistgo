# im2dhistgo

## Introduction
`im2dhistgo` is a Go package designed for image processing. It provides functions to generate 1D and 2D histograms.

### 2D histogram
A moving window of WxW moves through out the given image, and as its center places on each pixel, number of encounters with same and other brightness intensities is counted seperately.
![How moving window works](https://raw.githubusercontent.com/Mamdasn/im2dhist/main/assets/how-it-works-window-kernel-title.jpg "How moving window works")
## Installation

To install `im2dhistgo`, use the following `go get` command:

```bash
go get github.com/yourusername/im2dhistgo
```

## Usage

```go
package main

import (
	"fmt"
	"os"
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

```
