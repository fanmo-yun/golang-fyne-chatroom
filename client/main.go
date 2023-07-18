package main

import (
	"log"
	"net"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var data = binding.BindStringList(
	&[]string{},
)

func handle(client net.Conn) {
	for {
		readbuffer := [1024]byte{}
		n, err := client.Read(readbuffer[:])
		if err != nil {
			log.Println(err)
			break
		}
		data.Append(string(readbuffer[:n]))
	}

}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Golang Chatroom")

	client, err := net.Dial("tcp", "127.0.0.1:8085")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer client.Close()

	go handle(client)

	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	input := widget.NewMultiLineEntry()

	bottom := container.New(layout.NewVBoxLayout(), widget.NewSeparator(), container.New(layout.NewAdaptiveGridLayout(2), input, widget.NewButton("send", func() {
		if _, err := client.Write([]byte(input.Text)); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	})))
	c := container.New(layout.NewBorderLayout(nil, bottom, nil, nil), bottom, list)

	myWindow.SetContent(c)

	myWindow.Resize(fyne.Size{Width: 600, Height: 400})
	myWindow.ShowAndRun()

}
