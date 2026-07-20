package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
)

var captchaFont = map[rune][7]string{
	'0': {"111", "101", "101", "101", "101", "101", "111"},
	'1': {"010", "110", "010", "010", "010", "010", "111"},
	'2': {"111", "001", "001", "111", "100", "100", "111"},
	'3': {"111", "001", "001", "111", "001", "001", "111"},
	'4': {"101", "101", "101", "111", "001", "001", "001"},
	'5': {"111", "100", "100", "111", "001", "001", "111"},
	'6': {"111", "100", "100", "111", "101", "101", "111"},
	'7': {"111", "001", "001", "010", "010", "010", "010"},
	'8': {"111", "101", "101", "111", "101", "101", "111"},
	'9': {"111", "101", "101", "111", "001", "001", "111"},
	'A': {"010", "101", "101", "111", "101", "101", "101"},
	'B': {"110", "101", "101", "110", "101", "101", "110"},
	'C': {"111", "100", "100", "100", "100", "100", "111"},
	'D': {"110", "101", "101", "101", "101", "101", "110"},
	'E': {"111", "100", "100", "111", "100", "100", "111"},
	'F': {"111", "100", "100", "111", "100", "100", "100"},
	'G': {"111", "100", "100", "101", "101", "101", "111"},
	'H': {"101", "101", "101", "111", "101", "101", "101"},
	'J': {"001", "001", "001", "001", "101", "101", "111"},
	'K': {"101", "101", "110", "100", "110", "101", "101"},
	'L': {"100", "100", "100", "100", "100", "100", "111"},
	'M': {"101", "111", "111", "101", "101", "101", "101"},
	'N': {"101", "111", "111", "111", "101", "101", "101"},
	'P': {"111", "101", "101", "111", "100", "100", "100"},
	'Q': {"111", "101", "101", "101", "111", "001", "001"},
	'R': {"111", "101", "101", "111", "110", "101", "101"},
	'T': {"111", "010", "010", "010", "010", "010", "010"},
	'U': {"101", "101", "101", "101", "101", "101", "111"},
	'W': {"101", "101", "101", "101", "111", "111", "101"},
	'X': {"101", "101", "101", "010", "101", "101", "101"},
	'Y': {"101", "101", "101", "010", "010", "010", "010"},
}

func GenerateCaptchaCode() string {
	const chars = "0123456789"
	code := make([]byte, 4)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}

func GenerateCaptchaImage(code string) (string, error) {
	img := image.NewRGBA(image.Rect(0, 0, 120, 40))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{R: 246, G: 248, B: 252, A: 255}}, image.Point{}, draw.Src)
	for i := 0; i < 90; i++ {
		img.Set(rand.Intn(120), rand.Intn(40), color.RGBA{R: uint8(90 + rand.Intn(80)), G: uint8(90 + rand.Intn(80)), B: uint8(90 + rand.Intn(80)), A: 255})
	}
	for i := 0; i < 4; i++ {
		drawChar(img, rune(code[i]), 10+i*27, 6, color.RGBA{R: uint8(rand.Intn(35)), G: uint8(rand.Intn(45)), B: uint8(40 + rand.Intn(70)), A: 255})
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func drawChar(img *image.RGBA, ch rune, x, y int, c color.Color) {
	pattern, ok := captchaFont[ch]
	if !ok {
		return
	}
	cell := 4
	for row, line := range pattern {
		for col, val := range line {
			if val != '1' {
				continue
			}
			rect := image.Rect(x+col*cell, y+row*cell, x+(col+1)*cell-1, y+(row+1)*cell-1)
			draw.Draw(img, rect, &image.Uniform{C: c}, image.Point{}, draw.Src)
		}
	}
}
