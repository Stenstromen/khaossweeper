package minesweeper

import (
	"fmt"
	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

type CustomToolbarItem struct {
	object fyne.CanvasObject
}

func (t *CustomToolbarItem) ToolbarObject() fyne.CanvasObject {
	return t.object
}

func NewCustomToolbarItem(object fyne.CanvasObject) widget.ToolbarItem {
	return &CustomToolbarItem{object: object}
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

	imgResource, err := fyne.LoadResourceFromPath("graphics/square.png")
	if err != nil {
		return fmt.Errorf("failed to load image: %v", err)
	}

	openedImgResource, err := fyne.LoadResourceFromPath("graphics/openedsquare.png")
	if err != nil {
		return fmt.Errorf("failed to load opened square image: %v", err)
	}

	mineImgResource, err := fyne.LoadResourceFromPath("graphics/mine.png")
	if err != nil {
		return fmt.Errorf("failed to load mine image: %v", err)
	}

	var resetGame func()

	myApp := app.NewWithID("com.github.stenstromen.khaossweeper")
	myWindow := myApp.NewWindow("KhaosSweeper")

	myWindow.Resize(fyne.NewSize(666, 666))
	myWindow.SetFixedSize(true)
	myWindow.SetIcon(mineImgResource)
	myWindow.SetTitle("KhaosSweeper")
	myWindow.CenterOnScreen()

	namespaceEntry := widget.NewEntry()
	namespaceEntry.Wrapping = fyne.TextWrapOff
	namespaceEntry.Resize(fyne.NewSize(300, namespaceEntry.MinSize().Height))

	podEntry := widget.NewEntry()
	podEntry.Wrapping = fyne.TextWrapOff
	podEntry.Resize(fyne.NewSize(250, podEntry.MinSize().Height))

	nsItem := container.NewHBox(
		widget.NewLabel(fmt.Sprint("Namespace: ", namespace)),
	)

	var podLabel *widget.Label = widget.NewLabel("")
	pItem := container.NewHBox(
		widget.NewLabel("Killed:"),
		podLabel,
	)

	currentNamespace := NewCustomToolbarItem(nsItem)
	currentPod := NewCustomToolbarItem(pItem)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderNewIcon(), func() {
			resetGame()
		}),
		widget.NewToolbarSpacer(),
		currentNamespace,
		currentPod,
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			dialog.ShowInformation("About", fmt.Sprintf("KhaosSweeper\n\nBy GitHub/Stenstromen\n\nSafe Mode: %v", safemode), myWindow)
		}),
	)

	const gridSize = 16
	const bombProbability = 0.1

	buttons := make([]fyne.CanvasObject, gridSize*gridSize)
	bombs := make([]bool, gridSize*gridSize)

	seed := int64(12345)
	myRand := rand.New(rand.NewSource(seed))

	for i := range bombs {
		if myRand.Float64() < bombProbability {
			bombs[i] = true
		}
	}

	updateGrid := func() {
		grid := container.NewGridWithColumns(gridSize)
		for i, btn := range buttons {
			if btn != nil {
				grid.Add(container.NewStack(btn))
			} else if bombs[i] {
				mineIcon := widget.NewIcon(mineImgResource)
				grid.Add(container.NewStack(mineIcon))
			} else {
				grid.Add(widget.NewLabel(""))
			}
		}
		myWindow.SetContent(container.NewBorder(toolbar, nil, nil, nil, grid))
	}

	removeSurroundingButtons := func(buttonID int) {
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if rand.Float32() < 0.5 {
					x := buttonID%gridSize + dx
					y := buttonID/gridSize + dy
					neighborID := y*gridSize + x
					if x >= 0 && x < gridSize && y >= 0 && y < gridSize && neighborID >= 0 && neighborID < len(buttons) {
						if bombs[neighborID] {
							_, podName := kubekiller.Kubekiller(kubeconfig, namespace, safemode)
							podLabel.SetText(podName)
							buttons[neighborID] = nil
							return
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
				_, podName := kubekiller.Kubekiller(kubeconfig, namespace, safemode)
				podLabel.SetText(podName)
				buttons[buttonID] = nil
				return
			} else {
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
						_, podName := kubekiller.Kubekiller(kubeconfig, namespace, safemode)
						podLabel.SetText(podName)
						buttons[buttonID] = nil
						return
					} else {
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
