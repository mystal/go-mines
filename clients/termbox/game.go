package game

import (
    "github.com/nsf/termbox-go"
    "github.com/Mystal/go-mines/minegrid"

    "fmt"
    "image"
)

type Difficulty int

const (
    DiffEasy Difficulty = iota
    DiffMedium
    DiffHard
)

var (
    grid *minegrid.MineGrid
    gridPosition image.Point
    cursorPos image.Point
    gridChanged bool
)

func Play() {
    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()

    gridPosition = image.Point{20, 0}
    initGame(DiffEasy)

    quit := false
    for !quit {
        display()
        quit = updateGame()
    }
}

func updateGame() bool {
    event := termbox.PollEvent()
    if event.Type == termbox.EventKey {
        if event.Ch == 'q' {
            return true
        }
        if event.Ch == 'n' {
            initGame(DiffEasy)
        } else if event.Key == termbox.KeySpace {
            gameState, _ := grid.Reveal(cursorPos.X, cursorPos.Y)
            if gameState != minegrid.GameContinue {
                drawGrid()
                termbox.HideCursor()
                if gameState == minegrid.GameWon {
                    drawMessage("You won! Press any key to exit.")
                } else {
                    drawMessage("You lost... Press any key to exit.")
                }
                termbox.Flush()
                termbox.PollEvent()
                return true
            }
            gridChanged = true
        } else if event.Ch == 'f' {
            grid.ToggleFlag(cursorPos.X, cursorPos.Y)
            gridChanged = true
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

func initGame(diff Difficulty) {
    if diff == DiffEasy {
        grid, _ = minegrid.MakeMineGrid(9, 9, 10)
    } else if diff == DiffMedium {
        grid, _ = minegrid.MakeMineGrid(16, 16, 40)
    } else if diff == DiffHard {
        grid, _ = minegrid.MakeMineGrid(40, 16, 99)
    }
    gridChanged = true

    cursorPos = image.Point{0, 0}

    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func moveUp() {
    if cursorPos.Y > 0 {
        cursorPos.Y -= 1
    }
}

func moveDown() {
    if cursorPos.Y < grid.Y() - 1 {
        cursorPos.Y += 1
    }
}

func moveLeft() {
    if cursorPos.X > 0 {
        cursorPos.X -= 1
    }
}

func moveRight() {
    if cursorPos.X < grid.X() - 1 {
        cursorPos.X += 1
    }
}

func display() {
    drawStatus()
    if gridChanged {
        drawGrid()
        gridChanged = false
    }
    termbox.SetCursor(cursorPos.X + gridPosition.X + 1, cursorPos.Y + gridPosition.Y + 1)
    termbox.Flush()
}

func drawGrid() {
    drawCells(colorGrid(grid.String()), gridPosition.X, gridPosition.Y)
}

func drawStatus() {
    drawString("Minesweeper", 0, 0) //TODO bold
    drawString(fmt.Sprintf("Mines: %02d", grid.MinesLeft()), 0, 3)
    drawActions()
}

func drawActions() {
    drawString("Actions", 0, 5) //TODO bold
    drawString("Space: reveal", 0, 6)
    drawString("f: toggle flag", 0, 7)
    drawString("n: new game", 0, 9)
    drawString("q: quit", 0, 10)
}

func drawMessage(str string) {
    drawString(str, 0, grid.Y() + 3)
}

func colorGrid(gridStr string) []termbox.Cell {
    cells := make([]termbox.Cell, len(gridStr))
    for i, c := range gridStr {
        switch c {
        case '-':
            cells[i] = termbox.Cell{' ', termbox.ColorWhite, termbox.ColorBlue}
        case 'F':
            cells[i] = termbox.Cell{c, termbox.ColorRed|termbox.AttrBold, termbox.ColorBlue}
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
            cells[i] = termbox.Cell{c, termbox.ColorWhite, termbox.ColorCyan}
        case '8':
            cells[i] = termbox.Cell{c, termbox.ColorWhite, termbox.ColorMagenta}
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
