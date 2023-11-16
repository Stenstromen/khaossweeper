package minesweeper

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/cobra"
	"github.com/stenstromen/khaossweeper/kubekiller"
)

func Minesweeper(cmd *cobra.Command, args []string) error {
	kubeconfig, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		return err
	}
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	}

	var resetGame func()

	myApp := app.New()
	myWindow := myApp.NewWindow("KhaosSweeper")

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			log.Println("New Game")
			resetGame()
		}),
		widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)

	/* 	content := container.NewBorder(toolbar, nil, nil, nil, widget.NewLabel("Content"))
	   	myWindow.SetContent(content) */
	myWindow.SetTitle("KhaosSweeper")

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
		//myWindow.SetContent(grid)
		myWindow.SetContent(container.NewBorder(toolbar, nil, nil, nil, grid))
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
							kubekiller.Kubekiller(kubeconfig, namespace)
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
				kubekiller.Kubekiller(kubeconfig, namespace)
			}
			buttons[buttonID] = nil
			removeSurroundingButtons(buttonID)
			updateGrid()
		})
	}

	resetGame = func() {
		for i := range bombs {
			bombs[i] = rand.Float64() < bombProbability
			buttonID := i // capture loop variable
			buttons[i] = widget.NewButton("X", func() {
				if bombs[buttonID] {
					fmt.Println("Boom! Hit a bomb!")
					kubekiller.Kubekiller(kubeconfig, namespace)
				}
				buttons[buttonID] = nil
				removeSurroundingButtons(buttonID)
				updateGrid()
			})
		}
		updateGrid()
	}

	//updateGrid()
	resetGame()

	myWindow.ShowAndRun()

	return nil
}
