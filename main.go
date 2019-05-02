package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var timer, x, y, currentType, currentRotate int

const (
	screenWidth = 240
	screenHeight = 320
	blockSize = 24
	maxBlockX = 10
	maxBlockY = 12
)

// block rendering
// create block
// move block

const (
	TYPE_1 = iota
	TYPE_2
	TYPE_3
	TYPE_4
	TYPE_5
	TYPE_6
	TYPE_7
)

const (
	ROTATE_0 = iota
	ROTATE_90
	ROTATE_180
	ROTATE_270
)

var blockTypes = map[int][][]int{
	TYPE_1: {
		{ 0, 0, 0, 0},
		{ 1, 1, 1, 1},
		{ 0, 0, 0, 0},
		{ 0, 0, 0, 0},
	},
	TYPE_2: {
		{1, 1},
		{1, 1},
	},
	TYPE_3: {
		{ 0, 1, 1},
		{ 1, 1, 0},
		{ 0, 0, 0},
	},
	TYPE_4: {
		{ 1, 1, 0},
		{ 0, 1, 1},
		{ 0, 0, 0},
	},
	TYPE_5: {
		{ 0, 1, 0},
		{ 1, 1, 1},
		{ 0, 0, 0},
	},
	TYPE_6: {
		{ 1, 0, 0},
		{ 1, 1, 1},
		{ 0, 0, 0},
	},
	TYPE_7: {
		{ 0, 0, 1},
		{ 1, 1, 1},
		{ 0, 0, 0},
	},
}

func update(screen *ebiten.Image) error {
	timer++
	if timer % 7 == 0 {
		handleInput()
	}
	ebitenutil.DebugPrint(screen, string(x))
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	//ebitenutil.DebugPrint(screen, "Hello, World!")
	//ebitenutil.DebugPrint(screen, "Hello, World!")
	draw(screen)
	return nil
}

var lines = make([][]int, 13)

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Hello, World!"); err != nil {
		log.Fatal(err)
	}
}

func handleInput() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if x > 0 {
			x--
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if x < maxBlockX {
			x++
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if y < maxBlockY {
			y++
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if currentRotate == ROTATE_270 {
			currentRotate = ROTATE_0
		} else {
			currentRotate++
		}
	}
}

func draw(screen *ebiten.Image) {
	for i, line := range lines {
		for j, block := range line {
			if block == 1 {
				img, _ := ebiten.NewImage(blockSize, blockSize, 0)
				img.Fill(color.RGBA{0xff, 0, 0, 0xff})
				options := &ebiten.DrawImageOptions{}
				options.GeoM.Translate(float64(j * 24), float64(i * 24))
				screen.DrawImage(img, options)
			}
		}
	}
	max := len(blockTypes[currentType])
	blocks := make([][]int, max)
	for i, _ := range blocks {
		blocks[i] = make([]int, max)
	}
	switch currentRotate {
	case ROTATE_0:
		for i, line := range blockTypes[currentType] {
			for j, block := range line {
				blocks[i][j] = block
			}
		}
	case ROTATE_90:
		for i, line := range blockTypes[currentType] {
			for j, block := range line {
				blocks[j][max-i-1] = block
			}
		}
	case ROTATE_180:
		for i, line := range blockTypes[currentType] {
			for j, block := range line {
				blocks[max-i-1][max-j-1] = block
			}
		}
	case ROTATE_270:
		for i, line := range blockTypes[currentType] {
			for j, block := range line {
				blocks[max-j-1][i] = block
			}
		}
	}
	for i, line := range blocks {
		for j, block := range line {
			if block == 1 {
				img, _ := ebiten.NewImage(blockSize, blockSize, 0)
				img.Fill(color.RGBA{120, 120, 120, 0xff})
				options := &ebiten.DrawImageOptions{}
				options.GeoM.Translate(float64((x + j) * blockSize), float64((y + i) * blockSize))
				screen.DrawImage(img, options)
			}
		}
	}
}

func init() {
	timer = 0
	x = 0
	y = 0
	currentType = TYPE_1
}


