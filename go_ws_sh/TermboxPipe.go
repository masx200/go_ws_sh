package go_ws_sh

import (
	"fmt"
	// "io"
	"log"

	"github.com/nsf/termbox-go"
)

func TermboxPipe(writable func(p []byte) (n int, err error), closable func() error) (onCancel func() error, startable func(), err error) {
	const ESCCH = 0x1B
	err = termbox.Init()
	if err != nil {
		log.Printf("termbox initialization failed: %v", err)
		return nil, nil, err
	}

	startable = func() {
		defer termbox.Close()
		termbox.SetCursor(0, 0)
		// 主循环
		defer func() { go closable() }()
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyTab:
					go writable([]byte{'\t'})
				// https://learn.microsoft.com/zh-cn/windows/console/console-virtual-terminal-sequences
				case termbox.KeySpace:
					// fmt.Println("Space key pressed")
					go writable([]byte{' '})
				case termbox.KeyArrowUp:
					go writable([]byte{0x1B, '[', 'A'})
				case termbox.KeyArrowDown:
					go writable([]byte{0x1B, '[', 'B'})
				case termbox.KeyArrowLeft:
					go writable([]byte{0x1B, '[', 'D'})
				case termbox.KeyArrowRight:
					go writable([]byte{0x1B, '[', 'C'})
				case termbox.KeyF1:
					go writable([]byte{ESCCH, 'O', 'P'})
				case termbox.KeyEnter:
					// fmt.Println("Enter key pressed")
					go writable([]byte{'\n'})
				case termbox.KeyBackspace:
					go writable([]byte{'\b'})
				// 	fmt.Println("Backspace key pressed")
				case termbox.KeyDelete:
					go writable([]byte{0x1B, '[', '3', '~'})
				// 	fmt.Println("Delete key pressed")
				case termbox.KeyHome:
					go writable([]byte{0x1B, '[', 'H'})
				// 	fmt.Println("Home key pressed")
				case termbox.KeyEnd:
					go writable([]byte{0x1B, '[', 'F'})
				// 	fmt.Println("End key pressed")
				case termbox.KeyEsc:
					go writable([]byte{0x1B})
				case termbox.KeyInsert:
					go writable([]byte{0x1B, '[', '2', '~'})
				case termbox.KeyPgup:
					go writable([]byte{0x1B, '[', '5', '~'})
				case termbox.KeyPgdn:
					go writable([]byte{0x1B, '[', '6', '~'})
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
					go writable([]byte{ESCCH, 'O', 'Q'})
				case termbox.KeyF3:
					go writable([]byte{ESCCH, 'O', 'R'})
				case termbox.KeyF4:
					go writable([]byte{ESCCH, 'O', 'S'})
				case termbox.KeyF5:
					go writable([]byte{ESCCH, '[', '2', '5', '~'})
				case termbox.KeyF6:
					go writable([]byte{ESCCH, '[', '1', '7', '~'})
				case termbox.KeyF7:
					go writable([]byte{ESCCH, '[', '1', '8', '~'})
				case termbox.KeyF8:
					go writable([]byte{ESCCH, '[', '1', '9', '~'})
				case termbox.KeyF9:
					go writable([]byte{ESCCH, '[', '2', '0', '~'})
				case termbox.KeyF10:
					go writable([]byte{ESCCH, '[', '2', '1', '~'})
				case termbox.KeyF11:
					go writable([]byte{ESCCH, '[', '2', '3', '~'})
				case termbox.KeyF12:
					go writable([]byte{ESCCH, '[', '2', '4', '~'})
					/* 无法理解为什么KeyCtrlSpace会导致翻车 */
				// case termbox.KeyCtrlSpace:
				// 	go writable([]byte{0x00})

				default:
					if ev.Ch != 0 {
						// fmt.Printf("Character '%c' (code: %d) was pressed\n", ev.Ch, ev.Ch)

						go writable([]byte{byte(ev.Ch)})
					} else if ev.Key < 256 {
						fmt.Printf("key event ascii with code: %d\n", ev.Key)
						//go writable([]byte{byte(ev.Key)})
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
	return func() error {
		termbox.Interrupt()
		return nil
	}, startable, nil
}
