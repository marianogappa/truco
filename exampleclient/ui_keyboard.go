//go:build !tinygo
// +build !tinygo

package exampleclient

import (
	"log"
	"os"

	"github.com/nsf/termbox-go"
)

func (u *ui) startKeyEventLoop() chan rune {
	keyPressesCh := make(chan rune)
	go func() {
		for {
			event := termbox.PollEvent()
			if event.Type != termbox.EventKey {
				continue
			}
			if event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC || event.Key == termbox.KeyCtrlD || event.Key == termbox.KeyCtrlZ || event.Ch == 'q' {
				termbox.Close()
				log.Println("Chau!")
				os.Exit(0)
			}
			keyPressesCh <- event.Ch
		}
	}()
	return keyPressesCh
}
