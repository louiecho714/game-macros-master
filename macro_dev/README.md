# Macro Dev

這個資料夾是獨立開發區，不會動到你原本的 `main.go`。
目前版本直接使用 Windows API，不依賴額外第三方套件。

## 功能

- 支援多組 `Macro` 設定
- 支援 `Ctrl + 1/2/3` 快速切換巨集
- 支援 `Ctrl + F1` 開始
- 支援 `Ctrl + F2` 暫停
- 支援鍵盤按鍵
- 支援滑鼠左鍵與右鍵
- 每個動作可分開設定按住時間與間隔時間
- `RepeatCount = 0` 代表無限循環

## 執行

在專案根目錄執行：

```powershell
go run .\macro_dev
```

## 熱鍵

- `Ctrl + 1`: 切換到第 1 組巨集
- `Ctrl + 2`: 切換到第 2 組巨集
- `Ctrl + 3`: 切換到第 3 組巨集
- `Ctrl + F1`: 執行目前選中的巨集
- `Ctrl + F2`: 暫停目前巨集

如果切換時剛好正在執行，程式會先停止目前巨集，再切換到新的設定。

## 設定方式

編輯 `macro_dev/main.go` 裡的 `macros := []Macro{ ... }`。

## 欄位說明

- `Hold`: 按鈕按下後持續多久才放開，單位是毫秒
- `Interval`: 這個動作做完後，等多久再做下一個動作，單位是毫秒
- `LoopDelayMS`: 每一輪跑完後，等多久再開始下一輪，單位是毫秒
- `RepeatCount`: 巨集執行次數，`0` 代表無限循環，直到你按 `Ctrl + F2`

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

## 可用動作

- `ActionKey`
- `ActionLeftClick`
- `ActionRightClick`

## 可用按鍵

目前內建：

- `0` 到 `9`
- `a` 到 `z`
- `space`
- `enter`
- `tab`
- `esc`
- `f1` 到 `f12`

如果你需要更多按鍵，我可以再幫你補進去。
