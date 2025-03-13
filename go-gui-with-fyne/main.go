package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

var data = [][]string{[]string{"top left", "top right"}, []string{"bottom left", "bottom right"}}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Table Widget")
	// resize
	myWindow.Resize(fyne.NewSize(400, 400))

	list := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		})

	myWindow.SetContent(list)
	myWindow.ShowAndRun()
}

// import (
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/layout"
// 	"fyne.io/fyne/v2/widget"
// )

// func main() {
// 	myApp := app.New()

// 	myWindow := myApp.NewWindow("Fyne Example")

// 	// Create a basic text label
// 	label := widget.NewLabel("Hello Fyne!")

// 	// Create a button with a callback function
// 	button := widget.NewButton("Quit", func() {
// 		myApp.Quit()
// 	})

// 	// Create a simple form with an entry field
// 	entry := widget.NewEntry()
// 	form := &widget.Form{
// 		OnSubmit: func() {
// 			label.SetText("Entered : " + entry.Text)
// 		},
// 	}
// 	form.Append("Entry : ", entry)

// 	// Combine widgets into a container
// 	content := container.NewVBox(
// 		label,
// 		form,
// 		button,
// 	)

// 	// Center the content on the screen
// 	myWindow.SetContent(container.New(layout.NewCenterLayout(), content))

// 	// Show the window
// 	myWindow.ShowAndRun()
// }
