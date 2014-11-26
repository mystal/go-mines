package main

import (
	"github.com/mystal/go-mines/minegrid"
	"github.com/nsf/termbox-go"

	"fmt"
	"image"
)

type (
	Difficulty   int
	gameState    int
	actionFunc   func() gameState
	actionLookup map[termbox.Event]actionFunc
)

const (
	DiffEasy Difficulty = iota
	DiffMedium
	DiffHard
)

const (
	statePlay gameState = iota
	stateLose
	stateWin
	stateNew
	stateQuit
)

var actionStrings = [...][]string{
	//statePlay
	[]string{
		"Space: reveal",
		"f: flag",
		"Arrow keys: move",
		"",
		"n: new game",
		"q: quit"},
	//stateLose
	[]string{
		"n: new game",
		"q: quit"},
	//stateWin
	[]string{
		"n: new game",
		"q: quit"},
	//stateNew
	[]string{
		"e: easy",
		"m: medium",
		"h: hard",
		"",
		"c: cancel",
		"q: quit"}}

var actionFuncs = [...]actionLookup{
	//statePlay
	actionLookup{
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeySpace}:      revealSquare,
		termbox.Event{Type: termbox.EventKey, Ch: 'f'}:                    flagSquare,
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeyArrowUp}:    moveUp,
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeyArrowDown}:  moveDown,
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeyArrowLeft}:  moveLeft,
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeyArrowRight}: moveRight,
		termbox.Event{Type: termbox.EventKey, Ch: 'n'}:                    newGame,
		termbox.Event{Type: termbox.EventKey, Ch: 'q'}:                    quitGame},
	//stateLose
	actionLookup{
		termbox.Event{Type: termbox.EventKey, Ch: 'n'}: newGame,
		termbox.Event{Type: termbox.EventKey, Ch: 'q'}: quitGame},
	//stateWin
	actionLookup{
		termbox.Event{Type: termbox.EventKey, Ch: 'n'}: newGame,
		termbox.Event{Type: termbox.EventKey, Ch: 'q'}: quitGame},
	//stateNew
	actionLookup{
		termbox.Event{Type: termbox.EventKey, Ch: 'e'}: newEasyGame,
		termbox.Event{Type: termbox.EventKey, Ch: 'm'}: newMediumGame,
		termbox.Event{Type: termbox.EventKey, Ch: 'h'}: newHardGame,
		termbox.Event{Type: termbox.EventKey, Ch: 'c'}: cancelNewGame,
		termbox.Event{Type: termbox.EventKey, Ch: 'q'}: quitGame}}

var (
	grid            *minegrid.MineGrid
	gridPosition    image.Point
	actionsPosition image.Point
	minesPosition   image.Point
	cursorPos       image.Point
	gridChanged     bool
)

func revealSquare() gameState {
	gridChanged, _ = grid.Reveal(cursorPos.X, cursorPos.Y)
	switch grid.State() {
	case minegrid.GridLost:
		return stateLose
	case minegrid.GridWon:
		return stateWin
	}
	return statePlay
}

func flagSquare() gameState {
	gridChanged, _ = grid.ToggleFlag(cursorPos.X, cursorPos.Y)
	return statePlay
}

func quitGame() gameState {
	return stateQuit
}

func newGame() gameState {
	return stateNew
}

func cancelNewGame() gameState {
	switch grid.State() {
	case minegrid.GridWon:
		return stateWin
	case minegrid.GridLost:
		return stateLose
	case minegrid.GridContinue:
		return statePlay
	}
	//TODO an error technically
	return stateQuit
}

func newEasyGame() gameState {
	initGame(DiffEasy)
	return statePlay
}

func newMediumGame() gameState {
	initGame(DiffMedium)
	return statePlay
}

func newHardGame() gameState {
	initGame(DiffHard)
	return statePlay
}

func initGame(diff Difficulty) {
	if diff == DiffEasy {
		grid, _ = minegrid.MakeMineGrid(9, 9, 10)
	} else if diff == DiffMedium {
		grid, _ = minegrid.MakeMineGrid(16, 16, 40)
	} else if diff == DiffHard {
		grid, _ = minegrid.MakeMineGrid(40, 16, 99)
	}

	minesPosition = image.Point{gridPosition.X + grid.X()/2, 0}
	cursorPos = image.Point{0, 0}
}

func moveUp() gameState {
	if cursorPos.Y > 0 {
		cursorPos.Y -= 1
	}
	return statePlay
}

func moveDown() gameState {
	if cursorPos.Y < grid.Y()-1 {
		cursorPos.Y += 1
	}
	return statePlay
}

func moveLeft() gameState {
	if cursorPos.X > 0 {
		cursorPos.X -= 1
	}
	return statePlay
}

func moveRight() gameState {
	if cursorPos.X < grid.X()-1 {
		cursorPos.X += 1
	}
	return statePlay
}

func display(curState gameState, clear bool) {
	if clear {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	}

	drawColorString("Minesweeper", 0, 0, termbox.AttrBold, termbox.ColorDefault)
	drawActions(curState)
	if gridChanged || clear {
		drawColorString(fmt.Sprintf("%02d", grid.MinesLeft()),
			minesPosition.X, minesPosition.Y,
			termbox.ColorRed|termbox.AttrBold, termbox.ColorWhite)
		drawCells(colorGrid(grid.String()), gridPosition.X, gridPosition.Y)
		gridChanged = false
	}
	drawStatus(curState)

	if curState == statePlay {
		termbox.SetCursor(cursorPos.X+gridPosition.X+1, cursorPos.Y+gridPosition.Y+1)
	} else {
		termbox.HideCursor()
	}
	termbox.Flush()
}

func drawActions(curState gameState) {
	curActionPos := actionsPosition.Y
	for _, action := range actionStrings[curState] {
		drawString(action, 0, curActionPos)
		curActionPos += 1
	}
}

func drawStatus(curState gameState) {
	status := ""
	switch curState {
	case statePlay:
		status = "Play!"
	case stateLose:
		status = "You lost..."
	case stateWin:
		status = "You won!"
	case stateNew:
		status = "Choose a difficulty"
	}
	drawString(status, 0, gridPosition.Y+grid.Y()+3)
}

func drawMessage(str string) {
	drawString(str, 0, grid.Y()+3)
}

func colorGrid(gridStr string) []termbox.Cell {
	cells := make([]termbox.Cell, len(gridStr))
	for i, c := range gridStr {
		switch c {
		case '-':
			cells[i] = termbox.Cell{' ', termbox.ColorWhite, termbox.ColorBlue}
		case 'F':
			cells[i] = termbox.Cell{c, termbox.ColorRed | termbox.AttrBold, termbox.ColorBlue}
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
	drawColorString(str, x, y, termbox.ColorDefault, termbox.ColorDefault)
}

func drawColorString(str string, x, y int, fg, bg termbox.Attribute) {
	i, j := x, y
	for _, c := range str {
		if c == '\n' {
			i = x
			j += 1
			continue
		}
		termbox.SetCell(i, j, c, fg, bg)
		i += 1
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	gridPosition = image.Point{20, 1}
	actionsPosition = image.Point{0, 2}

	initGame(DiffEasy)

	clear := true
	for curState := statePlay; curState != stateQuit; {
		display(curState, clear)
		clear = false
		action := actionFuncs[curState][termbox.PollEvent()]
		if action != nil {
			nextState := action()
			clear = nextState != curState
			curState = nextState
		}
	}
}
