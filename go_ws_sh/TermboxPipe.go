package go_ws_sh

import (
	"fmt"
	"io"
	"log"

	"github.com/nsf/termbox-go"
)

type TermboxCloser struct {
}

// Close implements io.Closer.
func (t *TermboxCloser) Close() error {
	termbox.Interrupt()
	return nil
}

func init() {
	var _ io.Closer = &TermboxCloser{}
}
func TermboxPipe(writable io.Writer, closable io.Closer) (onCancel io.Closer, startable func(), err error) {
	err = termbox.Init()
	if err != nil {
		log.Printf("termbox initialization failed: %v", err)
		return nil, nil, err
	}

	startable = func() {
		defer termbox.Close()
		termbox.SetCursor(0, 0)
		// 主循环
		defer closable.Close()
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				// case termbox.KeySpace:
				// 	fmt.Println("Space key pressed")
				// case termbox.KeyArrowUp:
				// 	fmt.Println("Up arrow key pressed")
				// case termbox.KeyArrowDown:
				// 	fmt.Println("Down arrow key pressed")
				// case termbox.KeyArrowLeft:
				// 	fmt.Println("Left arrow key pressed")
				// case termbox.KeyArrowRight:
				// 	fmt.Println("Right arrow key pressed")
				// case termbox.KeyEnter:
				// 	fmt.Println("Enter key pressed")
				// case termbox.KeyBackspace:
				// 	fmt.Println("Backspace key pressed")
				// case termbox.KeyDelete:
				// 	fmt.Println("Delete key pressed")
				// case termbox.KeyHome:
				// 	fmt.Println("Home key pressed")
				// case termbox.KeyEnd:
				// 	fmt.Println("End key pressed")
				// case termbox.KeyEsc:
				// 	fmt.Println("ESC key pressed")
				case termbox.KeyCtrlC:
					fmt.Println("CtrlC key pressed exit")
					closable.Close()
					return // 退出程序
				case termbox.KeyCtrlD:
					fmt.Println("CtrlD key pressed exit")
					closable.Close()
					return // 退出程序
				case termbox.KeyCtrlZ:
					fmt.Println("CtrlZ key pressed exit")
					closable.Close()
					return // 退出程序
				default:
					if ev.Ch != 0 {
						fmt.Printf("Character '%c' (code: %d) was pressed\n", ev.Ch, ev.Ch)

						writable.Write([]byte{byte(ev.Ch)})
					} else {
						fmt.Printf("Unknown key event with code: %d\n", ev.Key)
					}
				}
			case termbox.EventError:
				log.Printf("Error event: %v", ev.Err)
				closable.Close()
				return
			}
		}
	}
	return &TermboxCloser{}, startable, nil
}
