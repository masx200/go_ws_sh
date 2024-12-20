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
					writable([]byte{'\n'})
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
				default:
					if ev.Ch != 0 {
						// fmt.Printf("Character '%c' (code: %d) was pressed\n", ev.Ch, ev.Ch)

						writable([]byte{byte(ev.Ch)})
					} else if ev.Key < 256 {
						fmt.Printf("key event ascii with code: %d\n", ev.Key)
						//writable([]byte{byte(ev.Key)})
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
