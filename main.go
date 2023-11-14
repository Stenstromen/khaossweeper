package main

import (
	"fmt"
	"math/rand"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("KhaosSweeper")

	const gridSize = 16
	const bombProbability = 0.1 // 10% chance of a bomb

	buttons := make([]*widget.Button, gridSize*gridSize)
	bombs := make([]bool, gridSize*gridSize) // Track if a button is a bomb

	rand.Seed(time.Now().UnixNano()) // Initialize the random number generator

	// Randomly place bombs
	for i := range bombs {
		if rand.Float64() < bombProbability {
			bombs[i] = true
		}
	}

	updateGrid := func() {
		grid := container.NewGridWithColumns(gridSize)
		for i, btn := range buttons {
			if btn != nil {
				grid.Add(container.NewStack(btn))
			} else if bombs[i] {
				grid.Add(widget.NewLabel("💣")) // Show a bomb icon
			} else {
				grid.Add(widget.NewLabel(""))
			}
		}
		myWindow.SetContent(grid)
	}

	removeSurroundingButtons := func(buttonID int) {
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if rand.Float32() < 0.5 { // 50% chance to consider this neighbor
					x := buttonID%gridSize + dx
					y := buttonID/gridSize + dy
					neighborID := y*gridSize + x
					if x >= 0 && x < gridSize && y >= 0 && y < gridSize && neighborID >= 0 && neighborID < len(buttons) {
						if bombs[neighborID] {
							fmt.Println("Boom! Hit a bomb!")
						}
						buttons[neighborID] = nil
					}
				}
			}
		}
	}

	for i := 0; i < gridSize*gridSize; i++ {
		buttonID := i
		buttons[i] = widget.NewButton("X", func() {
			if bombs[buttonID] {
				fmt.Println("Boom! Hit a bomb!")
			}
			buttons[buttonID] = nil
			removeSurroundingButtons(buttonID)
			updateGrid()
		})
	}

	updateGrid()

	myWindow.ShowAndRun()
}
