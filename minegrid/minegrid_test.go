package minegrid

import (
	"testing"
)

func TestMakeMindGrid_Empty(t *testing.T) {
	width := 10
	height := 5
	mines := 0

	grid, err := MakeMineGrid(width, height, mines)

	if err != nil {
		t.Fatalf("Could not create grid: %s", err)
	}

	if grid.X() != width {
		t.Errorf("Width is wrong!\nExpected: %d\nActual: %d", width, grid.X())
	}
	if grid.Y() != height {
		t.Errorf("Height is wrong!\nExpected: %d\nActual: %d", height, grid.Y())
	}

	if grid.MinesLeft() != mines {
		t.Errorf("Mines is wrong!\nExpected: %d\nActual: %d", mines, grid.MinesLeft())
	}
}

func TestMakeMindGrid_TooManyMines(t *testing.T) {
	width := 10
	height := 5
	mines := 100

	grid, err := MakeMineGrid(width, height, mines)

	if grid != nil || err == nil {
		t.Fatalf("Should have failed to create grid!")
	}
}

func TestCheckPoint(t *testing.T) {
	grid, _ := MakeMineGrid(10, 10, 0)

	if err := grid.checkPoint(0, 0); err != nil {
		t.Errorf("%d, %d should be valid! Got: %s", err)
	}
	if err := grid.checkPoint(5, 5); err != nil {
		t.Errorf("%d, %d should be valid! Got: %s", err)
	}
	if err := grid.checkPoint(9, 9); err != nil {
		t.Errorf("%d, %d should be valid! Got: %s", err)
	}

	if err := grid.checkPoint(-1, 0); err == nil {
		t.Errorf("%d, %d should be valid! Got: %s", err)
	}
	if err := grid.checkPoint(0, -1); err == nil {
		t.Errorf("%d, %d should be valid! Got: %s", err)
	}
	if err := grid.checkPoint(10, 0); err == nil {
		t.Errorf("%d, %d should be valid! Got: %s", err)
	}
	if err := grid.checkPoint(0, 10); err == nil {
		t.Errorf("%d, %d should be valid! Got: %s", err)
	}
}
