package minegrid

import (
	"fmt"
	"image"
	"math/rand"
	"time"
)

type cell struct {
	point            image.Point
	mines            uint8
	flags            uint8
	revealed         bool
	surroundingMines uint8
}

type MineGrid struct {
	cells        [][]cell
	x, y         int
	mines        int
	maxFlags     uint8
	minesFlagged int
	spacesLeft   int
	state        GridState
}

type GridState uint

type Error string

const (
	GridContinue GridState = iota
	GridWon
	GridLost
)

func (e Error) Error() string {
	return string(e)
}

func MakeMineGrid(x, y, mines int) (*MineGrid, error) {
	if mines > x*y {
		return nil, Error(fmt.Sprintf("Too many mines for a %dx%d grid", x, y))
	}

	rand.Seed(time.Now().Unix())

	mineSet := make(map[image.Point]bool)
	for i := 0; i < mines; i++ {
		p := image.Point{rand.Intn(x), rand.Intn(y)}
		for mineSet[p] {
			p.X, p.Y = rand.Intn(x), rand.Intn(y)
		}
		mineSet[p] = true
	}

	cells := make([][]cell, y)
	for j := 0; j < y; j++ {
		cells[j] = make([]cell, x)
		for i := 0; i < x; i++ {
			p := image.Point{i, j}
			cells[j][i].point = p
			if mineSet[p] {
				cells[j][i].mines = 1
			}
		}
	}
	g := MineGrid{cells, x, y, mines, 1, 0, x*y - mines, GridContinue}
	//Count surrounding mines and store
	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			g.cells[j][i].surroundingMines, _ = g.countSurroundingMines(i, j)
		}
	}

	return &g, nil
}

func (g MineGrid) X() int {
	return g.x
}

func (g MineGrid) Y() int {
	return g.y
}

func (g MineGrid) State() GridState {
	return g.state
}

func (g MineGrid) MinesLeft() int {
	return g.mines - g.minesFlagged
}

func (g MineGrid) String() (str string) {
	if g.x == 0 || g.y == 0 {
		return
	}

	str += "#"
	for i := 0; i < g.x; i++ {
		str += "#"
	}
	str += "#\n"
	for j := 0; j < g.y; j++ {
		str += "#"
		for i := 0; i < g.x; i++ {
			if g.cells[j][i].flags != 0 {
				str += "F"
			} else if !g.cells[j][i].revealed {
				str += "-"
			} else if g.cells[j][i].mines != 0 {
				str += "*"
			} else if g.cells[j][i].surroundingMines != 0 {
				str += fmt.Sprint(g.cells[j][i].surroundingMines)
			} else {
				str += " "
			}
		}
		str += "#\n"
	}
	str += "#"
	for i := 0; i < g.x; i++ {
		str += "#"
	}
	str += "#\n"
	return
}

func (g MineGrid) checkPoint(x, y int) error {
	if y < 0 {
		return Error("y is negative")
	}
	if y >= g.y {
		return Error("y is too big")
	}
	if x < 0 {
		return Error("x is negative")
	}
	if x >= g.x {
		return Error("x is too big")
	}
	return nil
}

func (g MineGrid) HasMine(x, y int) (bool, error) {
	if err := g.checkPoint(x, y); err != nil {
		return false, err
	}

	return g.cells[y][x].mines != 0, nil
}

func (g MineGrid) GetNeighbors(x, y int) ([]image.Point, error) {
	if err := g.checkPoint(x, y); err != nil {
		return nil, err
	}

	neighbors := make([]image.Point, 0, 8)
	for j := y - 1; j <= y+1; j++ {
		for i := x - 1; i <= x+1; i++ {
			if (x != i || y != j) && g.checkPoint(i, j) == nil {
				neighbors = append(neighbors, image.Point{i, j})
			}
		}
	}
	return neighbors, nil
}

func (g MineGrid) countSurroundingMines(x, y int) (uint8, error) {
	if err := g.checkPoint(x, y); err != nil {
		return 0, err
	}

	points, _ := g.GetNeighbors(x, y)
	count := uint8(0)
	for _, p := range points {
		count += g.cells[p.Y][p.X].mines
	}
	return count, nil
}

func (g MineGrid) countSurroundingFlags(x, y int) (uint8, error) {
	if err := g.checkPoint(x, y); err != nil {
		return 0, err
	}

	points, _ := g.GetNeighbors(x, y)
	count := uint8(0)
	for _, p := range points {
		count += g.cells[p.Y][p.X].flags
	}
	return count, nil
}

func (g *MineGrid) ToggleFlag(x, y int) (bool, error) {
	if err := g.checkPoint(x, y); err != nil {
		return false, err
	}

	if g.cells[y][x].revealed {
		return false, nil
	}

	if g.cells[y][x].flags == g.maxFlags {
		g.minesFlagged -= int(g.maxFlags)
		g.cells[y][x].flags = 0
	} else {
		g.cells[y][x].flags += 1
		g.minesFlagged += 1
	}
	return true, nil
}

func (g *MineGrid) Reveal(x, y int) (bool, error) {
	if err := g.checkPoint(x, y); err != nil {
		return false, err
	}

	if g.cells[y][x].flags != 0 {
		return false, nil
	}

	if g.cells[y][x].revealed {
		if g.cells[y][x].surroundingMines == 0 {
			return false, nil
		}
		flags, _ := g.countSurroundingFlags(x, y)
		if g.cells[y][x].surroundingMines != flags {
			return false, nil
		}
		neighbors, _ := g.GetNeighbors(x, y)
		for _, p := range neighbors {
			if !g.cells[p.Y][p.X].revealed {
				g.Reveal(p.X, p.Y)
			}
		}
		return true, nil
	}

	g.cells[y][x].revealed = true
	if g.cells[y][x].mines != 0 {
		g.state = GridLost
		return true, nil
	}
	g.spacesLeft -= 1
	if g.spacesLeft == 0 {
		g.state = GridWon
		return true, nil
	}
	if g.cells[y][x].surroundingMines == 0 {
		neighbors, _ := g.GetNeighbors(x, y)
		for _, p := range neighbors {
			g.Reveal(p.X, p.Y)
		}
	}
	return true, nil
}
