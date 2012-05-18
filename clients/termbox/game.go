package game

import (
    "github.com/nsf/termbox-go"
    "mystal/minesweeper/minegrid"

    //"fmt"
    "image"
)

var grid minegrid.MineGrid
var gridChanged bool
var cursorPos image.Point

func Play() {
    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()

    //TODO initialize things
    grid, _ = minegrid.MakeMineGrid(9, 9, 10)
    gridChanged = true

    quit := false
    for !quit {
        drawGrid()
        quit = updateGame()
    }
}

func updateGame() bool {
    event := termbox.PollEvent()
    if event.Type == termbox.EventKey {
        if event.Ch == 'q' {
            return true
        }
        if event.Key == termbox.KeySpace {
            state, _ := grid.Reveal(cursorPos.X, cursorPos.Y)
            //drawString(fmt.Sprint(minegrid.SpacesLeft), len(grid[0]) + 3, 0)
            //termbox.Flush()
            if state != 0 {
                drawGrid()
                termbox.HideCursor()
                if state == minegrid.GameWon {
                    drawString("You won! Press any key to exit.", 0, len(grid) + 3)
                } else {
                    drawString("You lost... Press any key to exit.", 0, len(grid) + 3)
                }
                termbox.Flush()
                termbox.PollEvent()
                return true
            }
            gridChanged = true
        } else if event.Ch == 'f' {
            //TODO flag as a mine
            //TODO update on-screen counter of mines unflagged
        }

        switch event.Key {
        case termbox.KeyArrowUp:
            moveUp()
        case termbox.KeyArrowDown:
            moveDown()
        case termbox.KeyArrowLeft:
            moveLeft()
        case termbox.KeyArrowRight:
            moveRight()
        }
    }
    return false
}

func moveUp() {
    if cursorPos.Y > 0 {
        cursorPos.Y -= 1
    }
    gridChanged = true
}

func moveDown() {
    if cursorPos.Y < len(grid) - 1 {
        cursorPos.Y += 1
    }
    gridChanged = true
}

func moveLeft() {
    if cursorPos.X > 0 {
        cursorPos.X -= 1
    }
    gridChanged = true
}

func moveRight() {
    if cursorPos.X < len(grid[0]) - 1 {
        cursorPos.X += 1
    }
    gridChanged = true
}

func drawGrid() {
    if gridChanged {
        drawCells(colorGrid(grid.String()), 0, 0)
        termbox.SetCursor(cursorPos.X + 1, cursorPos.Y + 1)
        termbox.Flush()
        gridChanged = false
    }
    gridChanged = true
}

func colorGrid(gridStr string) []termbox.Cell {
    cells := make([]termbox.Cell, len(gridStr))
    for i, c := range gridStr {
        switch c {
        case '-':
            cells[i] = termbox.Cell{' ', termbox.ColorWhite, termbox.ColorBlue}
        case '*':
            cells[i] = termbox.Cell{c, termbox.ColorWhite, termbox.ColorRed}
        case '1':
            cells[i] = termbox.Cell{c, termbox.ColorBlue, termbox.ColorDefault}
        case '2':
            cells[i] = termbox.Cell{c, termbox.ColorGreen, termbox.ColorDefault}
        case '3':
            cells[i] = termbox.Cell{c, termbox.ColorRed, termbox.ColorDefault}
        case '4':
            cells[i] = termbox.Cell{c, termbox.ColorYellow, termbox.ColorDefault}
        case '5':
            cells[i] = termbox.Cell{c, termbox.ColorMagenta, termbox.ColorDefault}
        case '6':
            cells[i] = termbox.Cell{c, termbox.ColorCyan, termbox.ColorDefault}
        case '7':
            cells[i] = termbox.Cell{c, termbox.ColorWhite, termbox.ColorBlue}
        case '8':
            cells[i] = termbox.Cell{c, termbox.ColorWhite, termbox.ColorGreen}
        default:
            cells[i] = termbox.Cell{c, termbox.ColorDefault, termbox.ColorDefault}
        }
    }
    return cells
}

func drawCells(cells []termbox.Cell, x, y int) {
    i, j := x, y
    for _, c := range cells {
        if c.Ch == '\n' {
            i = x
            j += 1
            continue
        }
        termbox.SetCell(i, j, c.Ch, c.Fg, c.Bg)
        i += 1
    }
}

func drawString(str string, x, y int) {
    i, j := x, y
    for _, c := range str {
        if c == '\n' {
            i = x
            j += 1
            continue
        }
        termbox.SetCell(i, j, c, termbox.ColorDefault, termbox.ColorDefault)
        i += 1
    }
}
