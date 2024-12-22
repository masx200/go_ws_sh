package go_ws_sh

import (
	"fmt"
	// "io"
	"log"

	"github.com/nsf/termbox-go"
)

func TermboxPipe(writable func(p []byte) (n int, err error), closable func() error, onsizechange func(cols int, rows int)) (onCancel func() error, startable func(), cols int, rows int, err error) {
	const ESCCH = 0x1B
	err = termbox.Init()
	if err != nil {
		log.Printf("termbox initialization failed: %v", err)
		return nil, nil, 0, 0, err
	}
	cols, rows = termbox.Size()
	startable = func() {
		defer termbox.Close()
		termbox.SetCursor(0, 0)
		// 主循环
		defer func() { go closable() }()
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventResize:
				cols, rows = ev.Width, ev.Height
				onsizechange(cols, rows)
			// case termbox.EventRaw:
			// 	log.Println(
			// 		"raw event: ", ev)
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyTab:
					/* 这里不能开协程会乱序不可以 */
					writable([]byte{'\t'})
				// https://learn.microsoft.com/zh-cn/windows/console/console-virtual-terminal-sequences
				case termbox.KeySpace:
					// fmt.Println("Space key pressed")
					writable([]byte{' '})
				case termbox.KeyArrowUp:
					writable([]byte{0x1B, '[', 'A'})
				case termbox.KeyArrowDown:
					writable([]byte{0x1B, '[', 'B'})
				case termbox.KeyArrowLeft:
					writable([]byte{0x1B, '[', 'D'})
				case termbox.KeyArrowRight:
					writable([]byte{0x1B, '[', 'C'})
				case termbox.KeyF1:
					writable([]byte{ESCCH, 'O', 'P'})
				case termbox.KeyEnter:
					// fmt.Println("Enter key pressed")
					writable([]byte{'\r'})
				case termbox.KeyBackspace:
					writable([]byte{'\b'})
				// 	fmt.Println("Backspace key pressed")
				case termbox.KeyDelete:
					writable([]byte{0x1B, '[', '3', '~'})
				// 	fmt.Println("Delete key pressed")
				case termbox.KeyHome:
					writable([]byte{0x1B, '[', 'H'})
				// 	fmt.Println("Home key pressed")
				case termbox.KeyEnd:
					writable([]byte{0x1B, '[', 'F'})
				// 	fmt.Println("End key pressed")
				case termbox.KeyEsc:
					writable([]byte{0x1B})
				case termbox.KeyInsert:
					writable([]byte{0x1B, '[', '2', '~'})
				case termbox.KeyPgup:
					writable([]byte{0x1B, '[', '5', '~'})
				case termbox.KeyPgdn:
					writable([]byte{0x1B, '[', '6', '~'})
				case termbox.KeyCtrlC:
					fmt.Println("CtrlC key pressed exit")
					go closable()
					return // 退出程序
				case termbox.KeyCtrlD:
					fmt.Println("CtrlD key pressed exit")
					go closable()
					return // 退出程序
				case termbox.KeyCtrlZ:
					fmt.Println("CtrlZ key pressed exit")
					go closable()
					return // 退出程序
				case termbox.KeyF2:
					writable([]byte{ESCCH, 'O', 'Q'})
				case termbox.KeyF3:
					writable([]byte{ESCCH, 'O', 'R'})
				case termbox.KeyF4:
					writable([]byte{ESCCH, 'O', 'S'})
				case termbox.KeyF5:
					writable([]byte{ESCCH, '[', '2', '5', '~'})
				case termbox.KeyF6:
					writable([]byte{ESCCH, '[', '1', '7', '~'})
				case termbox.KeyF7:
					writable([]byte{ESCCH, '[', '1', '8', '~'})
				case termbox.KeyF8:
					writable([]byte{ESCCH, '[', '1', '9', '~'})
				case termbox.KeyF9:
					writable([]byte{ESCCH, '[', '2', '0', '~'})
				case termbox.KeyF10:
					writable([]byte{ESCCH, '[', '2', '1', '~'})
				case termbox.KeyF11:
					writable([]byte{ESCCH, '[', '2', '3', '~'})
				case termbox.KeyF12:
					writable([]byte{ESCCH, '[', '2', '4', '~'})
					/* 无法理解为什么KeyCtrlSpace会导致翻车 */
				// case termbox.KeyCtrlSpace:
				// 	 writable([]byte{0x00})

				default:
					/* Ctrl 的传递方式通常与从系统接收完全相同。 这通常是向下移到控制字符保留空间 (0x0-
					0x1f) 的单个字符。 例 */
					if ev.Key <= termbox.KeyCtrl8 && termbox.KeyCtrlA <= ev.Key {
						writable([]byte{byte(ev.Key)})
					} else if ev.Ch != 0 {
						// fmt.Printf("Character '%c' (code: %d) was pressed\n", ev.Ch, ev.Ch)

						writable([]byte{byte(ev.Ch)})
					} else if ev.Key < 256 {
						fmt.Printf("key event ascii with code: %d\n", ev.Key)
						// writable([]byte{byte(ev.Key)})
					} else {
						fmt.Printf("key event unknown with code: %d\n", ev.Key)

					}
				}
			case termbox.EventError:
				log.Printf("Error event: %v", ev.Err)
				go closable()
				return
			}
		}
	}
	x := func() error {
		termbox.Interrupt()
		return nil
	}
	return x, startable, cols, rows, nil
}
