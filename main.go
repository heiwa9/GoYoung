//go:generate fyne bundle -o bundle.go -append myapp.png
//go:generate fyne bundle -o bundle.go -append AlibabaPuHuiTi-2-55-Regular.ttf
//
//fyne package -os linux -icon myapp.png
//fyne package -os windows -icon myapp.png
//fyne package -os android -appID cn.corehub.goyoung -icon myapp.png
package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	version = "Version 1.1.0"
	myApp   MyApp
	message = "Hello World!\n"
	user    User
)

type MyApp struct {
	app      fyne.App
	win      fyne.Window
	lab      *widget.Label
	msg      *widget.Entry
	prefix   *widget.Select
	username *widget.Entry
	password *widget.Entry
}

func main() {
	myApp = MyApp{app: app.NewWithID("GoYoung")}
	myApp.win = myApp.app.NewWindow("GoYoung")
	myApp.win.Resize(fyne.Size{Width: 320, Height: 480})

	myApp.lab = widget.NewLabel(version)
	myApp.msg = widget.NewMultiLineEntry()
	myApp.prefix = widget.NewSelect([]string{"!^Adcm0", "!^Iqnd0", "!^Mswx0"}, func(value string) {})
	myApp.prefix.SetSelectedIndex(0)
	myApp.username = widget.NewEntry()
	myApp.username.SetPlaceHolder("Account")
	myApp.password = widget.NewPasswordEntry()
	myApp.password.SetPlaceHolder("Password")
	user.ReadUserInfoJson()
	myApp.prefix.SetSelected(user.UserHard)
	myApp.username.SetText(user.UserAccount)
	myApp.password.SetText(user.PassWord)
	message += user.CheckServer(baiduURL) + "\n"
	myApp.msg.SetText(message)

	form := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items: []*widget.FormItem{
			{Text: "Account", Widget: myApp.username},
			{Text: "Password", Widget: myApp.password},
		},
		OnSubmit:   submitHandle,
		OnCancel:   cancelHandle,
		SubmitText: "Login",
		CancelText: "Offline",
	}

	content := container.New(layout.NewGridLayoutWithColumns(1), myApp.msg,
		container.NewVBox(myApp.lab, myApp.prefix, form))
	myApp.win.SetContent(content)
	myApp.win.ShowAndRun()
}

func submitHandle() {
	user.UserHard = myApp.prefix.Selected
	user.UserAccount = strings.TrimSpace(myApp.username.Text)
	user.PassWord = strings.TrimSpace(myApp.password.Text)
	loginInfo := user.login()
	if loginInfo == "50：认证成功" {
		user.SaveUserInfoJson()
		message += loginInfo + "\n"
		myApp.msg.SetText(message)
		dialog.ShowInformation("Login...", "LOGIN SUCCESSFUL", myApp.win)
		return
	}
	dialog.ShowInformation("Error", "LOGIN FAILED", myApp.win)
}

func cancelHandle() {
	if user.LastLoginURL != "" {
		logoutInfo := user.Logout()
		message += logoutInfo + "\n"
		myApp.msg.SetText(message)
		dialog.ShowInformation("Offline...", logoutInfo, myApp.win)
		return
	}
	message += "Please login with Go Young first!\n"
	myApp.msg.SetText(message)
	dialog.ShowInformation("Error", "Please login with Go Young first!", myApp.win)
}
