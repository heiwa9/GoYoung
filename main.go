//go:generate fyne package -os linux -icon myapp.png
//go:generate fyne package -os windows -icon myapp.png
//go:generate fyne package -os android/arm64 -appID cn.corehub.goyoung -icon myapp.png
package main

import (
	"GoYoung/lib"
	"errors"
	"os"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/qifengzhang007/goCurl"
)

const (
	version = "Version 1.3.1"
)

var (
	httpClient = goCurl.CreateHttpClient()
	filepath   = ""
)

func main() {
	var err error
	filepath, err = os.UserHomeDir()
	if err != nil {
		switch runtime.GOOS {
		case "android":
			filepath = "/storage/emulated/0/Android/data/cc.geekland.goyoung"
		default:
			filepath = "."
		}
	}
	filepath += "/.config/goyoung"
	_ = os.MkdirAll(filepath, os.ModePerm)
	filepath += "/user.json"

	myApp := app.NewWithID("GoYoung")
	myApp.Settings().SetTheme(&myTheme{})
	win := myApp.NewWindow("GoYoung")

	// 设置窗口的尺寸
	x, h := lib.ScreenSize()
	sx := float32(x)
	sh := float32(h)
	// 如果当前设备是移动端
	if fyne.CurrentDevice().IsMobile() {
		// 设置窗口的尺寸为最大
		win.Resize(fyne.NewSize(sx, sh))
	} else {
		// 设置窗口的尺寸为40% sy,(1-float32(sx)/float32(sh))*sh
		win.Resize(fyne.NewSize(400, 600))
	}

	// 设置窗口内容
	win.SetContent(pageHome(win))
	// 设置窗口居中
	win.CenterOnScreen()
	// 启动并展示窗口
	win.ShowAndRun()
}

func pageHome(win fyne.Window) fyne.CanvasObject {
	prefix := widget.NewSelect([]string{"!^Adcm0", "!^Iqnd0", "!^Mswx0"}, func(value string) {})
	prefix.SetSelectedIndex(0)
	username := widget.NewEntry()
	username.SetPlaceHolder("账号")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("密码")
	var user User
	err := lib.ReadJsonBind(filepath, &user)
	if err != nil {
		dialog.ShowError(err, win)
	}
	prefix.SetSelected(user.UserHard)
	username.SetText(user.UserAccount)
	password.SetText(user.PassWord)

	message := lib.CheckServer(baiduURL) + "\n"

	msg := widget.NewMultiLineEntry()
	msg.SetText(message)

	form := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items: []*widget.FormItem{
			{Text: "账号", Widget: username},
			{Text: "密码", Widget: password},
		},
		OnSubmit: func() {
			user.UserHard = prefix.Selected
			user.UserAccount = strings.TrimSpace(username.Text)
			user.PassWord = strings.TrimSpace(password.Text)
			loginInfo := user.login()
			if loginInfo == "50：认证成功" {
				err := lib.WriteJsonFile(filepath, &user)
				if err != nil {
					dialog.ShowError(err, win)
				}
				msg.SetText(message)
				dialog.ShowInformation("登录……", loginInfo, win)
				return
			}
			message += loginInfo + "\n"
			msg.SetText(message)
			dialog.ShowError(errors.New(loginInfo), win)
		},
		OnCancel: func() {
			if user.LastLoginURL != "" {
				logoutInfo := user.Logout()
				message += logoutInfo + "\n"
				msg.SetText(message)
				dialog.ShowInformation("下线……", logoutInfo, win)
				return
			}
			message += "请先使用GoYoung登录!\n"
			msg.SetText(message)
			dialog.ShowError(errors.New("请先使用GoYoung登录"), win)
		},
		SubmitText: "登录",
		CancelText: "下线",
	}
	return container.NewPadded(container.NewVSplit(msg, container.NewVBox(container.NewGridWithColumns(4, widget.NewLabel(runtime.GOOS), widget.NewLabel(version)), container.NewGridWithColumns(2, widget.NewLabel("账号前缀"), prefix), form)))
}
