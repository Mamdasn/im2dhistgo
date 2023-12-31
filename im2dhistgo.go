package im2dhistgo

import (
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/gif"
	_ "image/png"
	"math"
	"os"

	"fmt"
)

func Im2dhist_file(imagename string)  [65536]int {
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

func Im2dhist(input_layer *image.Gray, w int) [65536]int {
	img_bounds := input_layer.Bounds()
	var twodhist [65536]int

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

					twodhist[index1] += v_diff + 1
					if v_1 == v_2 {
						continue
					}
					index2 := uint16(v_2) + uint16(v_1)*256
					twodhist[index2] += v_diff + 1
				}
			}

		}
	}

	return twodhist
}

func Imhist(img *image.Gray) [256]int {
	var histogram [256]int
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			g, _, _, _ := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			histogram[g>>8]++
		}
	}
	return histogram
}

func float2uint8(c float64) uint8 {
	return uint8(math.Round(c * 255))
}

func max(arr [256]float64) float64 {
	maxm := arr[0]
	for i, _ := range arr {
		if maxm < arr[i] {
			maxm = arr[i]
		}
	}
	return maxm
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

type HSV struct { // {0..1}
	H, S, V float64
}

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
		// r = g = b = 0
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

type RGB struct { // {0..255}
	R, G, B float64
}

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
