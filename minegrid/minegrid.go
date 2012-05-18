package minegrid

import (
    "fmt"
    "image"
    "math/rand"
    "time"
)

type cell struct {
    point image.Point
    mines uint8
    flags uint8
    revealed bool
    sorroundingMines uint8
}

type MineGrid [][]cell

type GameState int

type Error string

const (
    GameContinue GameState = iota
    GameWon
    GameLost
)

var spacesLeft int

func (e Error) Error() string {
    return string(e)
}

func MakeMineGrid(x, y, mines int) (MineGrid, error) {
    if mines > x*y {
        return nil, Error(fmt.Sprintf("Too many mines for a %dx%d grid", x, y))
    }

    spacesLeft = x*y - mines

    rand.Seed(time.Now().Unix())

    mineSet := make(map[image.Point]bool)
    for i := 0; i < mines; i++ {
        p := image.Point{rand.Intn(x), rand.Intn(y)}
        for mineSet[p] {
            p.X, p.Y = rand.Intn(x), rand.Intn(y)
        }
        mineSet[p] = true
    }

    g := MineGrid(make([][]cell, y))
    for j := 0; j < y; j++ {
        g[j] = make([]cell, x)
        for i := 0; i < x; i++ {
            p := image.Point{i, j}
            g[j][i].point = p
            if mineSet[p] {
                g[j][i].mines = 1
            }
        }
    }
    //Count sorrounding mines and store
    for j := 0; j < y; j++ {
        for i := 0; i < x; i++ {
            g[j][i].sorroundingMines, _ = g.countSorroundingMines(i, j)
        }
    }

    return g, nil
}

func (g MineGrid) String() (str string) {
    if len(g) == 0 {
        return
    }

    str += "/"
    for i := 0; i < len(g[0]); i++ {
        str += "\u00AF"
    }
    str += "\\\n"
    for y := 0; y < len(g); y++ {
        str += "|"
        for x := 0; x < len(g[y]); x++ {
            if !g[y][x].revealed {
                str += "-"
            } else if g[y][x].mines != 0 {
                str += "*"
            } else if g[y][x].sorroundingMines != 0 {
                str += fmt.Sprint(g[y][x].sorroundingMines)
            } else {
                str += " "
            }
        }
        str += "|\n"
    }
    str += "\\"
    for i := 0; i < len(g[0]); i++ {
        str += "_"
    }
    str += "/\n"
    return
}

func (g MineGrid) checkPoint(x, y int) error {
    if y < 0 {
        return Error("y is negative")
    }
    if y >= len(g) {
        return Error("y is too big")
    }
    if x < 0 {
        return Error("x is negative")
    }
    if x >= len(g[0]) {
        return Error("x is too big")
    }
    return nil
}

func (g MineGrid) HasMine(x, y int) (bool, error) {
    if err := g.checkPoint(x, y); err != nil {
        return false, err
    }

    return g[y][x].mines != 0, nil
}

func (g MineGrid) GetNeighbors(x, y int) ([]image.Point, error) {
    if err := g.checkPoint(x, y); err != nil {
        return nil, err
    }

    neighbors := make([]image.Point, 0, 8)
    for j := y - 1; j <= y + 1; j++ {
        for i := x - 1; i <= x + 1; i++ {
            if (x != i || y != j) && g.checkPoint(i, j) == nil {
                neighbors = append(neighbors, image.Point{i, j})
            }
        }
    }
    return neighbors, nil
}

func (g MineGrid) countSorroundingMines(x, y int) (uint8, error) {
    if err := g.checkPoint(x, y); err != nil {
        return 0, err
    }

    points, _ := g.GetNeighbors(x, y)
    count := uint8(0)
    for i := 0; i < len(points); i++ {
        count += g[points[i].Y][points[i].X].mines
    }
    return count, nil
}

func (g MineGrid) Reveal(x, y int) (GameState, error) {
    if err := g.checkPoint(x, y); err != nil {
        return GameContinue, err
    }

    if g[y][x].revealed || g[y][x].flags != 0 {
        return GameContinue, nil
    }
    g[y][x].revealed = true
    if g[y][x].mines != 0 {
        return GameLost, nil
    }
    spacesLeft -= 1
    if spacesLeft == 0 {
        return GameWon, nil
    }
    if g[y][x].sorroundingMines == 0 {
        neighbors, _ := g.GetNeighbors(x, y)
        for _, p := range neighbors {
            g.Reveal(p.X, p.Y)
        }
    }
    return GameContinue, nil
}
