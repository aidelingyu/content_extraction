package main

import (
	"fmt"
	"fyne-gui/ocr"
	"log"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var selectFileType []string = []string{".pdf", ".jpg", ".jpeg", ".png"}

func mainShow(myWindow fyne.Window) {
	title := widget.NewLabel("")
	dirLabel := widget.NewLabel("目录文件:")
	filePath := widget.NewEntry() //文本输入框

	head := container.NewCenter(title)

	button1 := widget.NewButton("打开", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}

			filePath.SetText(reader.URI().String()) //把读取到的路径显示到输入框中
		}, myWindow)
		fd.SetFilter(storage.NewExtensionFileFilter(selectFileType))
		fd.Resize(fyne.NewSize(800, 500))
		fd.Show()
	})

	v1 := container.NewBorder(layout.NewSpacer(), layout.NewSpacer(), dirLabel, button1, filePath)

	text := widget.NewMultiLineEntry()

	filePathLabel := widget.NewLabel("")

	button2 := widget.NewButton("提取", func() {
		fp := filePath.Text
		if fp == "" {
			return
		}

		osSystem := runtime.GOOS

		path := ""

		switch osSystem {
		case "darwin":
			path = strings.Replace(fp, "file://", "", 1)
		case "windows":
			path = strings.Replace(fp, "file://", "", 1)
		}
		filePathLabel.SetText(fmt.Sprintf("操作系统：%s\n文件所在路径: %s", osSystem, path))

		resArray := ocr.GetTaxNumber(path)

		text.SetText(strings.Join(resArray, "\n"))
	})

	v2 := container.NewHBox(button2)
	v2Center := container.NewCenter(v2)
	ctnt := container.NewVBox(head, v1, v2Center, container.New(layout.NewGridWrapLayout(fyne.NewSize(990, 400)), text), filePathLabel)

	myWindow.SetContent(ctnt)
}
func main() {

	myApp := app.New()
	myApp.Settings().SetTheme(&myTheme{})
	myWindow := myApp.NewWindow("文字解析")
	myWindow.SetFixedSize(true)
	mainShow(myWindow)

	myWindow.Resize(fyne.NewSize(1000, 600))
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}
