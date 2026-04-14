# Macro Dev

`macro_dev` 是一個使用 Go 撰寫的 Windows 巨集工具，透過 Windows API 模擬鍵盤與滑鼠操作。

## 功能

- 支援多組 `Macro` 設定
- 支援 `Ctrl + 1/2/3` 切換巨集
- 支援 `Ctrl + F1` 啟動巨集
- 支援 `Ctrl + F2` 停止巨集
- 支援模擬鍵盤按鍵
- 支援滑鼠左鍵與右鍵點擊
- `RepeatCount = 0` 代表無限循環，直到手動停止

## 執行

在專案根目錄執行：

```powershell
go run .\macro_dev
```

## Build EXE

Windows 下可直接把 `macro_dev/main.go` 打包成 `.exe`：

```powershell
go build -o macro_dev.exe ./macro_dev
```

產生的檔案會在專案根目錄：

```text
game-macros-master\macro_dev.exe
```

## 操作方式

- `Ctrl + 1`：切換到第 1 組巨集
- `Ctrl + 2`：切換到第 2 組巨集
- `Ctrl + 3`：切換到第 3 組巨集
- `Ctrl + F1`：啟動目前選中的巨集
- `Ctrl + F2`：停止目前執行中的巨集

## 巨集設定

請編輯 `macro_dev/main.go` 裡的 `macros := []Macro{ ... }`。

## 參數說明

- `Hold`：按鍵或滑鼠按住的時間，單位為毫秒
- `Interval`：單一步驟執行後的等待時間，單位為毫秒
- `LoopDelayMS`：每輪巨集執行完後的等待時間
- `RepeatCount`：巨集重複次數，`0` 代表持續執行直到按下 `Ctrl + F2`

## 範例

```go
macros := []Macro{
    {
        Name:        "farm-basic",
        LoopDelayMS: 2000,
        RepeatCount: 0,
        Actions: []Action{
            {Type: ActionKey, Key: "1", Hold: 80, Interval: 500},
            {Type: ActionKey, Key: "2", Hold: 80, Interval: 500},
            {Type: ActionRightClick, Hold: 90, Interval: 500},
        },
    },
    {
        Name:        "boss-burst",
        LoopDelayMS: 1200,
        RepeatCount: 0,
        Actions: []Action{
            {Type: ActionKey, Key: "1", Hold: 90, Interval: 150},
            {Type: ActionKey, Key: "2", Hold: 90, Interval: 150},
            {Type: ActionKey, Key: "3", Hold: 90, Interval: 200},
        },
    },
}
```

## 支援的動作

- `ActionKey`
- `ActionLeftClick`
- `ActionRightClick`

## 支援的按鍵

- `0` 到 `9`
- `a` 到 `z`
- `space`
- `enter`
- `tab`
- `esc`
- `f1` 到 `f12`

如果要新增更多按鍵，可以在 `virtualKeyTable` 中擴充。
