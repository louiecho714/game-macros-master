package main

import (
	"fmt"
	"syscall"
	"time"
)

// Windows API 常量
const (
	KEYEVENTF_KEYUP = 0x0002 // 表示釋放按鍵

	MOUSEEVENTF_RIGHTDOWN = 0x0008 // 滑鼠右鍵按下
	MOUSEEVENTF_RIGHTUP   = 0x0010 // 滑鼠右鍵釋放

	MOUSEEVENTF_LEFTDOWN = 0x0002 // 滑鼠左鍵按下
	MOUSEEVENTF_LEFTUP   = 0x0004 // 滑鼠左鍵釋放

	VK_MBUTTON = 0x04 // 滑鼠中鍵
	VK_1       = 0x31 // 1 鍵
	VK_2       = 0x32 // 2 鍵
	VK_3       = 0x33 // 3 鍵
	VK_4       = 0x34 // 4 鍵
	VK_5       = 0x35 //5
	VK_SPACE   = 0x20 // 空白鍵 (Space)
	VK_F       = 0x46 // F 鍵
)
const (
	MOUSEEVENTF_MOVE     = 0x0001
	MOUSEEVENTF_ABSOLUTE = 0x8000
)

type POINT struct {
	X, Y int32
}

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procSetCursorPos     = user32.NewProc("SetCursorPos")
	procMouseEvent       = user32.NewProc("mouse_event")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

func getScreenSize() (width, height int32) {
	const (
		SM_CXSCREEN = 0
		SM_CYSCREEN = 1
	)
	w, _, _ := procGetSystemMetrics.Call(SM_CXSCREEN)
	h, _, _ := procGetSystemMetrics.Call(SM_CYSCREEN)
	return int32(w), int32(h)
}
func moveMouse(x, y int32) {
	procSetCursorPos.Call(uintptr(x), uintptr(y))
}

func clickMouse() {
	procMouseEvent.Call(MOUSEEVENTF_LEFTDOWN, 0, 0, 0, 0)
	time.Sleep(200 * time.Millisecond) // 模拟按下持续时间
	procMouseEvent.Call(MOUSEEVENTF_LEFTUP, 0, 0, 0, 0)
}

// 大俠立志傳 地圖移動
func funcMove() {
	// 获取屏幕尺寸
	screenWidth, screenHeight := getScreenSize()

	// 定义左上角和中央位置的坐标
	leftTopX, leftTopY := int32(0), int32(0)
	centerX, centerY := screenWidth, screenHeight

	// 移动到左上角并点击
	moveMouse(leftTopX, leftTopY)
	time.Sleep(1000 * time.Millisecond) // 等待鼠标移动完成
	clickMouse()

	// 等待一段时间
	time.Sleep(5 * time.Second)

	// 移动到中央并点击
	moveMouse(centerX, centerY)
	time.Sleep(1000 * time.Millisecond) // 等待鼠标移动完成
	clickMouse()

	// 等待一段时间
	time.Sleep(5 * time.Second)

}

// keybd_event 呼叫 Windows API，模擬按鍵事件
func keybd_event(bVk, bScan, dwFlags, dwExtraInfo uintptr) {
	mod := syscall.NewLazyDLL("user32.dll").NewProc("keybd_event")
	mod.Call(bVk, bScan, dwFlags, dwExtraInfo)
}

// mouse_event 呼叫 Windows API，模擬滑鼠事件
func mouse_event(dwFlags, dx, dy, dwData, dwExtraInfo uintptr) {
	mod := syscall.NewLazyDLL("user32.dll").NewProc("mouse_event")
	mod.Call(dwFlags, dx, dy, dwData, dwExtraInfo)
}

// pressKey 函數模擬按下和釋放按鍵
func pressKey(key uintptr) {
	keybd_event(key, 0, 0, 0)               // 按下鍵
	time.Sleep(100 * time.Millisecond)      // 模擬按下的時間
	keybd_event(key, 0, KEYEVENTF_KEYUP, 0) // 釋放鍵
}

// clickRightMouse 函數模擬滑鼠右鍵點擊
func clickRightMouse() {
	mouse_event(MOUSEEVENTF_RIGHTDOWN, 0, 0, 0, 0) // 按下右鍵
	time.Sleep(1000 * time.Millisecond)            // 按住一段時間
	mouse_event(MOUSEEVENTF_RIGHTUP, 0, 0, 0, 0)   // 釋放右鍵
}

// 模擬滑鼠左鍵點擊
func clickLeftMouse() {
	mouse_event(MOUSEEVENTF_LEFTDOWN, 0, 0, 0, 0) // 按下左鍵
	time.Sleep(100 * time.Millisecond)            // 按住一段時間
	mouse_event(MOUSEEVENTF_LEFTUP, 0, 0, 0, 0)   // 釋放左鍵
}

// 全域變數用來追蹤是否運行中
var isRunning = false
var spaceTime = 100 * time.Millisecond
var otherTime = 100 * time.Millisecond

// 快跑鴨
func function_1() {
	pressKey(VK_SPACE)
	time.Sleep(spaceTime)

	pressKey(VK_1)
	time.Sleep(otherTime)

	pressKey(VK_SPACE)
	time.Sleep(spaceTime)

	pressKey(VK_2)
	time.Sleep(otherTime)

	pressKey(VK_SPACE)
	time.Sleep(spaceTime)

	pressKey(VK_3)
	time.Sleep(otherTime)

	pressKey(VK_SPACE)
	time.Sleep(spaceTime)

	pressKey(VK_4)
	time.Sleep(otherTime)

	pressKey(VK_SPACE)
	time.Sleep(spaceTime)

	clickRightMouse() // 模擬右鍵點擊
	time.Sleep(otherTime)
}

// 一般腳色
func function_2() {
	clickLeftMouse()
	time.Sleep(spaceTime)
}

// 主要技能滑鼠右鍵連點
func function_3() {
	clickRightMouse()
	time.Sleep(spaceTime)

	pressKey(VK_1)
	time.Sleep(otherTime)

	clickRightMouse()
	time.Sleep(spaceTime)

	pressKey(VK_2)
	time.Sleep(otherTime)

	clickRightMouse()
	time.Sleep(spaceTime)

	pressKey(VK_3)
	time.Sleep(otherTime)

	clickRightMouse()
	time.Sleep(spaceTime)

	pressKey(VK_4)
	time.Sleep(otherTime)

}

// 只放1.2.3.4 輔助技能
func function_4() {

	pressKey(VK_1)
	time.Sleep(otherTime)

	pressKey(VK_2)
	time.Sleep(otherTime)

	pressKey(VK_3)
	time.Sleep(otherTime)

	pressKey(VK_4)
	time.Sleep(otherTime)

}

func function_5() {

	clickRightMouse()
	time.Sleep(spaceTime)

	pressKey(VK_1)
	time.Sleep(otherTime)

	clickRightMouse()
	time.Sleep(spaceTime)

	pressKey(VK_2)
	time.Sleep(otherTime)

	clickRightMouse()
	time.Sleep(spaceTime)

	pressKey(VK_3)
	time.Sleep(otherTime)

	clickRightMouse()
	time.Sleep(spaceTime)

	pressKey(VK_4)
	time.Sleep(otherTime)

	// clickRightMouse()
	// time.Sleep(spaceTime)

	// pressKey(VK_5)
	// time.Sleep(otherTime)

}

// 只放1 輔助技能
func function_only_one() {

	clickRightMouse()
	time.Sleep(otherTime)
	clickRightMouse()
	time.Sleep(otherTime)
	clickRightMouse()
	time.Sleep(5 * otherTime)
	pressKey(VK_F)
	time.Sleep(25 * otherTime)

}

// startLoop 控制無窮迴圈按鍵
func startLoop(stop chan bool) {

	for {
		select {
		case <-stop: // 接收到停止訊號，結束迴圈
			return
		default:
			function_only_one()
		}
	}
}

// listenForMiddleClick 監聽滑鼠中鍵事件，切換開始/停止
func listenForMiddleClick(stop chan bool) {
	user32 := syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState := user32.NewProc("GetAsyncKeyState")

	for {
		val, _, _ := getAsyncKeyState.Call(uintptr(VK_MBUTTON))
		if val != 0 { // 如果偵測到滑鼠中鍵按下
			if isRunning {
				fmt.Println("暫停...")
				stop <- true // 傳送停止訊號
				isRunning = false
			} else {
				fmt.Println("開始...")
				isRunning = true
				go startLoop(stop) // 開始新的無窮迴圈
			}
			time.Sleep(500 * time.Millisecond) // 避免重複偵測
		}
		time.Sleep(50 * time.Millisecond) // 減少 CPU 消耗
	}
}

func main() {
	stop := make(chan bool)       // 用來控制停止訊號的 channel
	go listenForMiddleClick(stop) // 開始監聽滑鼠中鍵

	// 防止主程式結束
	select {}
}
