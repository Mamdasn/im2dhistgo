// Package im2dhistgo provides functions for processing images and generating histograms.
package im2dhistgo

import (
	// Standard library imports
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/gif"
	_ "image/png"
	"math"
	"os"
	"sync"
	"fmt"
)

// Im2dhist_file takes an image file name as input, converts it to grayscale,
// and returns a 2D histogram of the image.
// The histogram is represented as a fixed-size array where each element
// counts occurrences of specific intensity differences.
func Im2dhist_file(imagename string)  [65536]uint32 {
	input_image, err := getImageFromFilePath(imagename)
	if err != nil {
		fmt.Println(imagename) // debugging
		fmt.Println("error:", err) // debugging
	}

	img_bounds := input_image.Bounds()
	m_v   := image.NewGray(img_bounds)

	// convert rgb image to hsv
	for x := img_bounds.Min.X; x < img_bounds.Max.X; x++ {
		for y := img_bounds.Min.Y; y < img_bounds.Max.Y; y++ {
			r, g, b, _ := input_image.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			rgb := &RGB{float64(r) / 65535, float64(g) / 65535, float64(b) / 65535}
			// convert rgb to hsv
			hsv := rgb.HSV()
			v := float2uint8(hsv.V)

			m_v.SetGray(x, y, color.Gray{Y: v})
		}
	}

	w := 1
	twodhist := Im2dhist(m_v, w)

	return twodhist
}

// Im2dhist generates a 2D histogram for a given grayscale image.
// The histogram is calculated based on the intensity differences of pixels
// within a specified window size.
func Im2dhist(input_layer *image.Gray, w int) [65536]uint32 {
	img_bounds := input_layer.Bounds()
	var twodhist [65536]uint32

	for x := img_bounds.Min.X; x < img_bounds.Max.X; x++ {
		for y := img_bounds.Min.Y; y < img_bounds.Max.Y; y++ {
			v_1 := input_layer.GrayAt(x, y).Y
			for i := -w; i < w+1; i++ {
				for j := -w; j < w+1; j++ {
					x_kernel := x + i
					y_kernel := y + j
					if x_kernel < img_bounds.Min.X || x_kernel >= img_bounds.Max.X || y_kernel < img_bounds.Min.Y || y_kernel >= img_bounds.Max.Y {
                    				continue
                			}
					v_2 := input_layer.GrayAt(x_kernel, y_kernel).Y

					index1 := uint16(v_1) + uint16(v_2)*256
					v_diff := int(v_2) - int(v_1)
					if v_diff < 0 {v_diff *= -1}

					v_diff_incremented := uint32(v_diff) + 1

					twodhist[index1] += v_diff_incremented
					if v_1 == v_2 {
						continue
					}
					index2 := uint16(v_2) + uint16(v_1)*256
					twodhist[index2] += v_diff_incremented
				}
			}

		}
	}

	return twodhist
}

// Im2dhist_parallel performs a similar operation to Im2dhist but uses
// parallel processing to speed up the histogram generation.
func Im2dhist_parallel(input_layer *image.Gray, w int) [65536]uint32 {
    img_bounds := input_layer.Bounds()
    var twodhist [65536]uint32
    var mutex sync.Mutex

    var wg sync.WaitGroup
    for x := img_bounds.Min.X; x < img_bounds.Max.X; x++ {
        wg.Add(1)
        go func(x int) {
            defer wg.Done()
            localHist := [65536]uint32{}
            for y := img_bounds.Min.Y; y < img_bounds.Max.Y; y++ {
                v_1 := input_layer.GrayAt(x, y).Y
                for i := -w; i < w+1; i++ {
                    for j := -w; j < w+1; j++ {
                        x_kernel := x + i
                        y_kernel := y + j
                        if x_kernel < img_bounds.Min.X || x_kernel >= img_bounds.Max.X || y_kernel < img_bounds.Min.Y || y_kernel >= img_bounds.Max.Y {
                            continue
                        }
                        v_2 := input_layer.GrayAt(x_kernel, y_kernel).Y

                        index1 := uint16(v_1) + uint16(v_2)*256
                        v_diff := int(v_2) - int(v_1)
                        if v_diff < 0 {
                            v_diff *= -1
                        }

                        v_diff_incremented := uint32(v_diff) + 1

                        localHist[index1] += v_diff_incremented
                        if v_1 != v_2 {
                            index2 := uint16(v_2) + uint16(v_1)*256
                            localHist[index2] += v_diff_incremented
                        }
                    }
                }
            }

            mutex.Lock()
            for i, v := range localHist {
                twodhist[i] += v
            }
            mutex.Unlock()
        }(x)
    }
    wg.Wait()

    return twodhist
}

// Imhist generates a standard histogram for a grayscale image.
// It returns an array where each element represents the count of pixels
// having a specific intensity value.
func Imhist(img *image.Gray) [256]uint32 {
	var histogram [256]uint32
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			v := img.GrayAt(x, y).Y
			histogram[uint8(v)]++
		}
	}
	return histogram
}

// float2uint8 converts a float64 value in the range [0, 1] to a uint8 value.
func float2uint8(c float64) uint8 {
	return uint8(math.Round(c * 255))
}

// max returns the maximum value in an array of float64.
func max(arr [256]float64) float64 {
	maxm := arr[0]
	for i, _ := range arr {
		if maxm < arr[i] {
			maxm = arr[i]
		}
	}
	return maxm
}

// getImageFromFilePath opens an image file and decodes it.
// It returns the decoded image and any error encountered.
func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

// HSV represents a color in the Hue, Saturation, and Value (HSV) color space.
type HSV struct { // {0..1}
	H, S, V float64
}

// HSV converts an RGB color to its HSV representation.
func (c *RGB) HSV() *HSV {
	r, g, b := c.R, c.G, c.B
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	v := max
	delta := max - min

	var h, s float64

	if max != 0 {
		s = delta / max
	} else {
		s = 0
		h = -1 // Undefined
		return &HSV{H: h, S: s, V: v}
	}

	if r == max {
		h = (g - b) / delta // Between yellow & magenta
	} else if g == max {
		h = 2 + (b-r)/delta // Between cyan & yellow
	} else {
		h = 4 + (r-g)/delta // Between magenta & cyan
	}

	h *= 60 // degrees
	if h < 0 {
		h += 360
	}

	// Normalize H to [0, 1)
	h = h / 360

	return &HSV{H: h, S: s, V: v}
}

// RGB represents a color in the Red, Green, and Blue (RGB) color space.
type RGB struct { // {0..255}
	R, G, B float64
}

// RGB converts an HSV color to its RGB representation.
func (c *HSV) RGB() *RGB {
	h, s, v := c.H, c.S, c.V

	h = math.Mod(h, 1.0) // Ensures h is within [0, 1)

	region := int(math.Floor(h * 6))
	fraction := h*6 - float64(region)
	p := v * (1.0 - s)
	q := v * (1 - fraction*s)
	t := v * (1 - (1-fraction)*s)

	var r, g, b float64

	switch region {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return &RGB{R: r, G: g, B: b}
}
