package minesweeper

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	//"fyne.io/fyne"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/cobra"
	"github.com/stenstromen/khaossweeper/kubekiller"
)

type ImageButton struct {
	widget.Icon
	OnTapped func()
	disabled bool
}

func NewImageButton(resource fyne.Resource, tapped func()) *ImageButton {
	img := &ImageButton{OnTapped: tapped}
	img.ExtendBaseWidget(img)
	img.SetResource(resource)
	return img
}

func (i *ImageButton) Tapped(*fyne.PointEvent) {
	if i.OnTapped != nil && !i.disabled {
		i.OnTapped()
	}
}

func (i *ImageButton) Disable() {
	i.disabled = true
	// Optional: Change the appearance to indicate the disabled state
	// i.SetResource(disabledImageResource)
}
func Minesweeper(cmd *cobra.Command, args []string) error {
	safemode, err := cmd.Flags().GetBool("safe-mode")
	if err != nil {
		return err
	}
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

	// Set the window size to 600x600
	myWindow.Resize(fyne.NewSize(600, 600))
	myWindow.SetFixedSize(true)

	imgResource, err := fyne.LoadResourceFromPath("graphics/square.png")
	if err != nil {
		return fmt.Errorf("failed to load image: %v", err)
	}

	openedImgResource, err := fyne.LoadResourceFromPath("graphics/openedsquare.png")
	if err != nil {
		return fmt.Errorf("failed to load opened square image: %v", err)
	}

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

	//buttons := make([]*widget.Button, gridSize*gridSize)

	buttons := make([]fyne.CanvasObject, gridSize*gridSize)
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
							kubekiller.Kubekiller(kubeconfig, namespace, safemode)
							buttons[neighborID] = nil
						} else {
							buttons[neighborID].(*ImageButton).Disable()
							buttons[neighborID].(*ImageButton).SetResource(openedImgResource)
						}
					}
				}
			}
		}
	}

	for i := 0; i < gridSize*gridSize; i++ {
		buttonID := i
		button := NewImageButton(imgResource, func() {
			if bombs[buttonID] {
				fmt.Println("Boom! Hit a bomb!")
				kubekiller.Kubekiller(kubeconfig, namespace, safemode)
			} else {
				// Directly update the button's image to the opened square
				buttons[buttonID].(*ImageButton).SetResource(openedImgResource)
			}
			buttons[buttonID].(*ImageButton).Disable()
			removeSurroundingButtons(buttonID)
			updateGrid()
		})

		buttons[i] = button
	}

	resetGame = func() {
		for i := range bombs {
			bombs[i] = rand.Float64() < bombProbability
			buttonID := i // capture loop variable

			button := NewImageButton(imgResource, func(buttonID int) func() {
				return func() {
					if bombs[buttonID] {
						fmt.Println("Boom! Hit a bomb!")
						kubekiller.Kubekiller(kubeconfig, namespace, safemode)
					} else {
						// Change the button's image to the opened square
						/* fyne.CurrentApp().Driver().CallOnMainThread(func() {
							buttons[buttonID].(*ImageButton).SetResource(openedImgResource)
						}) */
						buttons[buttonID].(*ImageButton).SetResource(openedImgResource)
					}
					buttons[buttonID].(*ImageButton).Disable()
					removeSurroundingButtons(buttonID)
					updateGrid()
				}
			}(buttonID))

			buttons[buttonID] = button
		}
		updateGrid()
	}

	resetGame()

	myWindow.ShowAndRun()

	return nil
}
