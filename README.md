# im2dhistgo

## Introduction
`im2dhistgo` is a Go package designed for image processing. It provides functions to generate 1D and 2D histograms.

### 2D histogram
A moving window of WxW moves through out the given image, and as its center places on each pixel, number of encounters with same and other brightness intensities is counted seperately.
![How moving window works](https://raw.githubusercontent.com/Mamdasn/im2dhist/main/assets/how-it-works-window-kernel-title.jpg "How moving window works")
## Installation

To install `im2dhistgo`, use the following `go get` command:

```bash
go get -u github.com/yourusername/im2dhistgo
```
