package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/k0kubun/pp"
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var timer, x, y, currentType, currentRotate, score, phase, nextType int

const (
	screenWidth = 320
	screenHeight = 320
	blockSize = 24
	maxBlockX = 10
	maxBlockY = 13
	initBlockX = 3
	scoreX = 200
	scoreY = 0
	nextX = 240
	nextY = 24
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

const (
	PHASE_GAMESTART = iota
	PHASE_GAMEOVER
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
	if phase != PHASE_GAMEOVER {
		timer++
		if timer % 7 == 0 {
			handleInput()
		}
	}
	if phase == PHASE_GAMEOVER {
		ebitenutil.DebugPrint(screen, "GAME OVER!: " + strconv.Itoa(score))
	} else {
		ebitenutil.DebugPrint(screen, strconv.Itoa(score))
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	//ebitenutil.DebugPrint(screen, "Hello, World!")
	//ebitenutil.DebugPrint(screen, "Hello, World!")
	draw(screen)
	return nil
}

var backgroundBlocks = make([][]int, maxBlockY)

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Hello, World!"); err != nil {
		log.Fatal(err)
	}
}

func handleInput() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if !isConflict(-1, 0) {
			x--
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if !isConflict(1, 0) {
			x++
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if !isConflict(0, 1) {
			y++
		} else {
			current := currentBlock()
			for i, line := range current {
				if y+i >= maxBlockY {
					break
				}
				for j, block := range line {
					if x+j >= maxBlockX {
						continue
					}
					if block == 1 {
						backgroundBlocks[y+i][x+j] = 1
					}
				}
			}
			handleLineClear()
			x = initBlockX
			y = 0
			currentRotate = ROTATE_0
			currentType = nextType
			nextType = generateBlock()
			if isConflict(0, 0) {
				phase = PHASE_GAMEOVER
			}
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		before := currentRotate
		if currentRotate == ROTATE_270 {
			currentRotate = ROTATE_0
		} else {
			currentRotate++
		}
		if isConflict(0, 0) {
			currentRotate = before
		}
	}
}

func draw(screen *ebiten.Image) {
	for i, line := range backgroundBlocks {
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
	current := currentBlock()
	for i, line := range current {
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
	for i, line := range blockTypes[nextType] {
		for j, block := range line {
			if block == 1 {
				img, _ := ebiten.NewImage(blockSize, blockSize, 0)
				img.Fill(color.RGBA{0x00, 0xff, 0x00, 0xff})
				options := &ebiten.DrawImageOptions{}
				options.GeoM.Translate(nextX + float64(j * blockSize), nextY + float64(i * blockSize))
				screen.DrawImage(img, options)
			}
		}
	}
	// text.Draw(screen, string(score), scoreFont, scoreX, scoreY, color.White)
}

func handleLineClear() {
	dy := 0
	newBackgroundBlocks := make([][]int, maxBlockY)
	for i := 0; i < maxBlockY; i++ {
		index := maxBlockY-i-1
		line := backgroundBlocks[index]
		newBackgroundBlocks[index] = make([]int, maxBlockX)
		if isLineClear(line) {
			dy++
			continue
		}
		for j, block := range line {
			newBackgroundBlocks[index+dy][j] = block
		}
	}
	switch dy {
	case 1:
		score += 100
	case 2:
		score += 300
	case 3:
		score += 600
	case 4:
		score += 1000
	}
	backgroundBlocks = newBackgroundBlocks
}

func isLineClear(line []int) bool {
	for _, block := range line {
		if block != 1 {
			return false
		}
	}
	return true
}

func currentBlock() [][]int {
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
	return blocks
}

func isConflict(dx int, dy int) bool {
	for i, line := range currentBlock() {
		newY := y+dy+i
		for j, block := range line {
			newX := x+dx+j
			if block == 1 {
				if newY >= maxBlockY || newX < 0 || newX >= maxBlockX || backgroundBlocks[newY][newX] == 1 {
					return true
				}
			}
		}
	}
	return false
}

func generateBlock() int {
	return rand.Intn(len(blockTypes))
}

func init() {
	timer = 0
	score = 0
	x = initBlockX
	y = 0
	rand.Seed(time.Now().UnixNano())
	currentType = generateBlock()
	nextType = generateBlock()
	phase = PHASE_GAMESTART
	for i, _ := range backgroundBlocks {
		backgroundBlocks[i] = make([]int, maxBlockX)
	}
}

func debug(args ...interface{}) {
	pp.Println(args...)
}


