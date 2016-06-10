package main

import (
	"github.com/andlabs/ui"
	"fmt"
	"regexp"
)

func main() {
	err := ui.Main(func() {
		dimention := ui.NewEntry()
		calculate := ui.NewButton("Calculate")
		box := ui.NewVerticalBox()
		matrix := ui.NewEntry()
		time := ui.NewLabel("23")
		result := ui.NewLabel("42")

		box.Append(ui.NewLabel("Dimention of the matrix or path to file:"), false)
		box.Append(dimention, false)
		box.Append(ui.NewLabel("Enter matrix values or leave blank for random:"), false)
		box.Append(matrix, false)
		box.Append(calculate, false)
		box.Append(ui.NewLabel("Time spent: "), false)
		box.Append(time, false)
		box.Append(ui.NewLabel("Result: "), false)
		box.Append(result, false)


		window := ui.NewWindow("Hello", 500, 300, false)
		window.SetChild(box)
		window.SetMargined(true)
		window.SetTitle("Determinant calculation")

		calculate.OnClicked(func(*ui.Button) {
			// greeting.SetText("Hello, " + name.Text() + "!")

			loadType := dimention.Text()
			match, _ := regexp.MatchString("([a-z]+)", loadType)
			if  !match {
				fmt.Println("Is Number")
				// call calc function with dimention and matrix
				// time.SetText()
				// result.SetText()
			} else {
				fmt.Println("Not Is Number")
				// load matrix from file with path dimentions, ignore matrix
				// time.SetText()
				// result.SetText()
			}
		})

		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}