package main

import (
	"log"
	"os"
	"fmt"
	"syscall"
	"php-proxy/icon"
	"github.com/getlantern/systray"
)

func main() {
	go func() {
		//
		log.SetOutput(os.Stdout)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		//
		config := &config{}
		config.init_config()
		//
		prx := &proxy{cfg: config}
		prx.init_proxy()
		//
		select {}
	}()
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("php-proxy")
	systray.SetTooltip("Windows Net Proxy")
	mShow := systray.AddMenuItem("显示", "显示窗口")
	mHide := systray.AddMenuItem("隐藏", "隐藏窗口")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出程序")

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	user32 := syscall.NewLazyDLL("user32.dll")
	getConsoleWindows := kernel32.NewProc("GetConsoleWindow")
	showWindowAsync := user32.NewProc("ShowWindowAsync")
	consoleHandle, r2, err := getConsoleWindows.Call()
	if consoleHandle == 0 {
		fmt.Println("Error call GetConsoleWindow: ", consoleHandle, r2, err)
	}

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				mShow.Disable()
				mHide.Enable()
				r1, r2, err := showWindowAsync.Call(consoleHandle, 5)
				if r1 != 1 {
					fmt.Println("Error call ShowWindow @SW_SHOW: ", r1, r2, err)
				}
			case <-mHide.ClickedCh:
				mHide.Disable()
				mShow.Enable()
				r1, r2, err := showWindowAsync.Call(consoleHandle, 0)
				if r1 != 1 {
					fmt.Println("Error call ShowWindow @SW_HIDE: ", r1, r2, err)
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
	// clean up here
}
