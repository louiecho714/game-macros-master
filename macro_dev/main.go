package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"syscall"
	"time"
)

type ActionType string

const (
	ActionKey        ActionType = "key"
	ActionLeftClick  ActionType = "left_click"
	ActionRightClick ActionType = "right_click"
)

const (
	keyEventKeyUp       = 0x0002
	mouseEventLeftDown  = 0x0002
	mouseEventLeftUp    = 0x0004
	mouseEventRightDown = 0x0008
	mouseEventRightUp   = 0x0010
	vkCtrl              = 0x11
	vkF1                = 0x70
	vkF2                = 0x71
	defaultHoldMS       = 60
)

type Action struct {
	Type     ActionType
	Key      string
	Hold     int
	Interval int
}

type Macro struct {
	Name        string
	LoopDelayMS int
	RepeatCount int
	Actions     []Action
}

var (
	user32            = syscall.NewLazyDLL("user32.dll")
	procKeybdEvent    = user32.NewProc("keybd_event")
	procMouseEvent    = user32.NewProc("mouse_event")
	procGetAsyncState = user32.NewProc("GetAsyncKeyState")
	virtualKeyTable   = map[string]uintptr{
		"1":     0x31,
		"2":     0x32,
		"3":     0x33,
		"4":     0x34,
		"5":     0x35,
		"6":     0x36,
		"7":     0x37,
		"8":     0x38,
		"9":     0x39,
		"0":     0x30,
		"a":     0x41,
		"b":     0x42,
		"c":     0x43,
		"d":     0x44,
		"e":     0x45,
		"f":     0x46,
		"g":     0x47,
		"h":     0x48,
		"i":     0x49,
		"j":     0x4A,
		"k":     0x4B,
		"l":     0x4C,
		"m":     0x4D,
		"n":     0x4E,
		"o":     0x4F,
		"p":     0x50,
		"q":     0x51,
		"r":     0x52,
		"s":     0x53,
		"t":     0x54,
		"u":     0x55,
		"v":     0x56,
		"w":     0x57,
		"x":     0x58,
		"y":     0x59,
		"z":     0x5A,
		"space": 0x20,
		"enter": 0x0D,
		"tab":   0x09,
		"esc":   0x1B,
		"f1":    0x70,
		"f2":    0x71,
		"f3":    0x72,
		"f4":    0x73,
		"f5":    0x74,
		"f6":    0x75,
		"f7":    0x76,
		"f8":    0x77,
		"f9":    0x78,
		"f10":   0x79,
		"f11":   0x7A,
		"f12":   0x7B,
	}
)

func main() {
	macros := []Macro{
		{
			Name:        "farm-basic",
			LoopDelayMS: 200,
			RepeatCount: 0,
			Actions: []Action{
				{Type: ActionKey, Key: "2", Hold: 90, Interval: 150},
				{Type: ActionKey, Key: "4", Hold: 90, Interval: 250},
				{Type: ActionRightClick, Hold: 80, Interval: 250},
				{Type: ActionKey, Key: "3", Hold: 90, Interval: 200},
				{Type: ActionKey, Key: "4", Hold: 90, Interval: 250},
				{Type: ActionRightClick, Hold: 80, Interval: 250},
			},
		},
		{
			Name:        "boss-burst",
			LoopDelayMS: 200,
			RepeatCount: 0,
			Actions: []Action{

				{Type: ActionKey, Key: "2", Hold: 90, Interval: 150},
				{Type: ActionKey, Key: "4", Hold: 90, Interval: 150},
				{Type: ActionRightClick, Hold: 80, Interval: 250},
				{Type: ActionKey, Key: "3", Hold: 90, Interval: 200},
				{Type: ActionKey, Key: "4", Hold: 90, Interval: 150},
				{Type: ActionRightClick, Hold: 80, Interval: 250},
			},
		},
	}

	runHotkeyLoop(macros)
}

func runHotkeyLoop(macros []Macro) {
	if len(macros) == 0 {
		log.Fatal("沒有可用的巨集設定")
	}

	var stop chan struct{}
	var isRunning atomic.Bool
	selectedIndex := 0

	printUsage(macros, selectedIndex)

	for {
		if number, ok := detectMacroSwitch(macros); ok {
			if isRunning.Load() {
				close(stop)
				isRunning.Store(false)
				fmt.Println("已停止目前巨集，準備切換設定")
			}

			selectedIndex = number
			fmt.Printf("已切換到巨集 [%d]: %s\n", selectedIndex+1, macros[selectedIndex].Name)
			waitHotkeyRelease(vkCtrl, keyNumberToVK(selectedIndex+1))
		}

		if isHotkeyPressed(vkCtrl, vkF1) && !isRunning.Load() {
			stop = make(chan struct{})
			isRunning.Store(true)
			currentMacro := macros[selectedIndex]
			fmt.Printf("巨集開始執行 [%d]: %s\n", selectedIndex+1, currentMacro.Name)

			go func(m Macro) {
				if err := runMacro(m, stop); err != nil {
					log.Printf("巨集執行失敗: %v\n", err)
				} else {
					log.Printf("巨集執行結束: %s\n", m.Name)
				}
				isRunning.Store(false)
			}(currentMacro)

			waitHotkeyRelease(vkCtrl, vkF1)
		}

		if isHotkeyPressed(vkCtrl, vkF2) && isRunning.Load() {
			close(stop)
			isRunning.Store(false)
			fmt.Println("巨集已暫停")
			waitHotkeyRelease(vkCtrl, vkF2)
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func printUsage(macros []Macro, selectedIndex int) {
	fmt.Println("Ctrl + F1 開始")
	fmt.Println("Ctrl + F2 暫停")
	fmt.Println("Ctrl + 1/2/3 切換巨集")
	for i, macro := range macros {
		marker := " "
		if i == selectedIndex {
			marker = "*"
		}
		fmt.Printf("[%s] %d: %s\n", marker, i+1, macro.Name)
	}
}

func detectMacroSwitch(macros []Macro) (int, bool) {
	maxSwitch := len(macros)
	if maxSwitch > 9 {
		maxSwitch = 9
	}

	for i := 0; i < maxSwitch; i++ {
		if isHotkeyPressed(vkCtrl, keyNumberToVK(i+1)) {
			return i, true
		}
	}

	return 0, false
}

func keyNumberToVK(number int) uintptr {
	return uintptr(0x30 + number)
}

func runMacro(m Macro, stop <-chan struct{}) error {
	if len(m.Actions) == 0 {
		return fmt.Errorf("沒有設定任何動作")
	}

	if m.RepeatCount == 0 {
		round := 1
		for {
			if shouldStop(stop) {
				return nil
			}

			if err := runRound(m, round, stop); err != nil {
				return err
			}
			round++
		}
	}

	for round := 1; round <= m.RepeatCount; round++ {
		if shouldStop(stop) {
			return nil
		}

		if err := runRound(m, round, stop); err != nil {
			return err
		}
	}

	return nil
}

func runRound(m Macro, round int, stop <-chan struct{}) error {
	fmt.Printf("[%s] 第 %d 輪開始\n", m.Name, round)

	for index, action := range m.Actions {
		if shouldStop(stop) {
			return nil
		}

		if err := runAction(action, stop); err != nil {
			return fmt.Errorf("第 %d 個動作失敗: %w", index+1, err)
		}
	}

	if m.LoopDelayMS > 0 {
		if err := sleepMSWithStop(m.LoopDelayMS, stop); err != nil {
			return nil
		}
	}

	return nil
}

func runAction(action Action, stop <-chan struct{}) error {
	holdMS := action.Hold
	if holdMS <= 0 {
		holdMS = defaultHoldMS
	}

	switch action.Type {
	case ActionKey:
		if action.Key == "" {
			return fmt.Errorf("key 動作缺少按鍵名稱")
		}
		fmt.Printf("按下鍵盤: %s, hold=%dms, interval=%dms\n", action.Key, holdMS, action.Interval)
		if err := tapKey(action.Key, holdMS); err != nil {
			return err
		}

	case ActionLeftClick:
		fmt.Printf("按下滑鼠左鍵: hold=%dms, interval=%dms\n", holdMS, action.Interval)
		leftClick(holdMS)

	case ActionRightClick:
		fmt.Printf("按下滑鼠右鍵: hold=%dms, interval=%dms\n", holdMS, action.Interval)
		rightClick(holdMS)

	default:
		return fmt.Errorf("不支援的動作類型: %s", action.Type)
	}

	if action.Interval > 0 {
		if err := sleepMSWithStop(action.Interval, stop); err != nil {
			return nil
		}
	}

	return nil
}

func tapKey(key string, holdMS int) error {
	vk, ok := virtualKeyTable[key]
	if !ok {
		return fmt.Errorf("不支援的按鍵: %s", key)
	}

	procKeybdEvent.Call(vk, 0, 0, 0)
	time.Sleep(ms(holdMS))
	procKeybdEvent.Call(vk, 0, keyEventKeyUp, 0)
	return nil
}

func leftClick(holdMS int) {
	procMouseEvent.Call(mouseEventLeftDown, 0, 0, 0, 0)
	time.Sleep(ms(holdMS))
	procMouseEvent.Call(mouseEventLeftUp, 0, 0, 0, 0)
}

func rightClick(holdMS int) {
	procMouseEvent.Call(mouseEventRightDown, 0, 0, 0, 0)
	time.Sleep(ms(holdMS))
	procMouseEvent.Call(mouseEventRightUp, 0, 0, 0, 0)
}

func isHotkeyPressed(modifier, key uintptr) bool {
	return isKeyDown(modifier) && isKeyDown(key)
}

func isKeyDown(key uintptr) bool {
	state, _, _ := procGetAsyncState.Call(key)
	return state&0x8000 != 0
}

func waitHotkeyRelease(keys ...uintptr) {
	for {
		allReleased := true
		for _, key := range keys {
			if isKeyDown(key) {
				allReleased = false
				break
			}
		}

		if allReleased {
			return
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func sleepMSWithStop(durationMS int, stop <-chan struct{}) error {
	timer := time.NewTimer(ms(durationMS))
	defer timer.Stop()

	select {
	case <-stop:
		return fmt.Errorf("stopped")
	case <-timer.C:
		return nil
	}
}

func shouldStop(stop <-chan struct{}) bool {
	select {
	case <-stop:
		return true
	default:
		return false
	}
}

func ms(value int) time.Duration {
	return time.Duration(value) * time.Millisecond
}
