package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	webview_go "github.com/webview/webview_go" // 注意这里改为webview_go
	"log"
	"path/filepath"
	"runtime"
)

func main() {
	myApp := app.New()
	window := myApp.NewWindow("测试环境")

	absPath, _ := filepath.Abs("assets/icon.png")
	icon, err := fyne.LoadResourceFromPath(absPath)
	if err != nil {
		log.Println("图标加载失败:", err)
	}

	//debug := runtime.GOOS != "windows" // Windows下建议关闭调试
	wv := webview_go.New(true)

	// 窗口设置
	wv.SetTitle("测试环境")
	// 设置窗口为全屏模式
	initWindow(window, wv)
	// 加载网页（可替换为你的H5地址）
	wv.Navigate("https://scrm-bq.test-chexiu.cn/login") // 示例页面

	//icon图标
	window.SetIcon(icon)
	myApp.SetIcon(icon)

	//// 可选：设置标题栏颜色
	window.SetTitle("北汽")
	window.SetPadded(false)
	window.CenterOnScreen()
	//// 窗口关闭时释放资源
	// 事件监控系统
	setupEventListeners(window, wv, myApp)
	defer func() {
		log.Println("执行延迟销毁")
		wv.Destroy()
	}()

	// 使用容器包装实现自动布局
	window.SetFullScreen(true)
	window.SetContent(widget.NewLabel("Loading...")) // 临时占位内容
	myApp.Run()
}

// 窗口关闭拦截器
func setupEventListeners(win fyne.Window, wv webview_go.WebView, myApp fyne.App) {
	// 1. 窗口关闭拦截器（主关闭事件）
	win.SetCloseIntercept(func() {
		log.Println("[事件] 用户触发关闭操作")

		// 执行资源释放
		log.Println("正在销毁WebView...")
		wv.Destroy()

		log.Println("关闭主窗口...")
		win.Close()

		log.Println("退出应用程序...")
		myApp.Quit()
	})

	// 2. 系统信号监控（仅Linux/macOS）
	if runtime.GOOS != "windows" {
	}

	// 3. WebView关闭事件（如果适用）
	wv.Bind("onClose", func() {
		log.Println("[WebView事件] 页面请求关闭")
		myApp.Quit()
	})
}

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

// Windows实现
//func getWindowsResolution() (int, int, error) {
//	user32 := syscall.NewLazyDLL("user32.dll")
//	getSystemMetrics := user32.NewProc("GetSystemMetrics")
//
//	cx, _, _ := getSystemMetrics.Call(SM_CXSCREEN)
//	if cx == 0 {
//		return 0, 0, fmt.Errorf("无法获取水平分辨率")
//	}
//
//	cy, _, _ := getSystemMetrics.Call(SM_CYSCREEN)
//	if cy == 0 {
//		return 0, 0, fmt.Errorf("无法获取垂直分辨率")
//	}
//
//	return int(cx), int(cy), nil
//}

// macOS实现
//func getMacResolution() (int, int, error) {
//	cmd := exec.Command("system_profiler", "SPDisplaysDataType")
//	output, err := cmd.CombinedOutput()
//	if err != nil {
//		return 0, 0, err
//	}
//
//	lines := strings.Split(string(output), "\n")
//	for _, line := range lines {
//		if strings.Contains(line, "Resolution:") {
//			parts := strings.Split(line, ":")
//			if len(parts) < 2 {
//				continue
//			}
//
//			res := strings.TrimSpace(parts[1])
//			res = strings.ReplaceAll(res, " ", "")
//			res = strings.ReplaceAll(res, "Retina", "")
//			res = strings.TrimSpace(res)
//
//			if dims := strings.Split(res, "x"); len(dims) == 2 {
//				width, _ := strconv.Atoi(dims[0])
//				height, _ := strconv.Atoi(dims[1])
//				return width, height, nil
//			}
//		}
//	}
//
//	return 0, 0, fmt.Errorf("未找到分辨率信息")
//}

func initWindow(window fyne.Window, wv webview_go.WebView) {
	// 获取主显示器尺寸（新API方式）
	width := 1024
	height := 768
	// 平台适配逻辑
	//if runtime.GOOS == "windows" {
	//	widths, heights, err := getWindowsResolution()
	//	if err == nil {
	//		width = widths
	//		height = heights
	//	}
	//} else if runtime.GOOS == "darwin" {
	//	widths, heights, err := getMacResolution()
	//	if err == nil {
	//		width = widths
	//		height = heights
	//	}
	//}
	fmt.Printf("Screen resolution: %dx%d\n", width, height)
	window.Resize(fyne.NewSize(1920, 1080-30)) // 任务栏补偿
	// 动态调整WebView尺寸（新版事件监听）
	wv.SetSize(width, height, webview_go.HintNone)
}
