package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	a := app.New()
	w := a.NewWindow("PGP Util")
	w.Resize(fyne.Size{
		Width:  512,
		Height: 512,
	})

	i := widget.NewEntry()
	i.SetPlaceHolder("Encrypted message")
	i.MultiLine = true

	oErr := widget.NewEntry()
	oErr.Disable()
	oErr.SetMinRowsVisible(2)
	oErr.MultiLine = true

	o := widget.NewEntry()
	o.SetMinRowsVisible(25)
	o.Disable()
	o.MultiLine = true

	content := container.NewVBox(
		i,
		widget.NewButton("Decrypt", func() {
			f, _ := os.CreateTemp(os.TempDir(), "pgp-gui-")
			f.WriteString(i.Text)
			cmd := exec.Command("gpg", "--decrypt", f.Name())
			rc, err := cmd.StderrPipe()
			sOut, _ := cmd.StdoutPipe()
			cmd.Start()
			log.Println(cmd.String())
			if err != nil {
				log.Fatal(err.Error())
			}
			stdErr, _ := io.ReadAll(rc)
			oErr.Text = string(stdErr)
			oErr.Refresh()
			stdOut, _ := io.ReadAll(sOut)

			o.Text = string(stdOut)
			o.Refresh()
		}),
		oErr,
		o,
	)
	w.SetContent(content)
	w.ShowAndRun()
}
