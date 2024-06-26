package exampleclient

import (
	"log"
	"os"
	"strconv"

	"github.com/nsf/termbox-go"
)

func (u *ui) startKeyEventLoop() {
	keyPressesCh := make(chan termbox.Event)
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
			keyPressesCh <- event
		}
	}()

	go func() {
		for {
			select {
			case <-keyPressesCh:
			case <-u.wantKeyPressCh:
				event := <-keyPressesCh
				u.sendKeyPressCh <- event.Ch
			}
		}
	}()
}

func (u *ui) pressAnyKey() {
	u.wantKeyPressCh <- struct{}{}
	<-u.sendKeyPressCh
}

func (u *ui) pressAnyNumber() int {
	u.wantKeyPressCh <- struct{}{}
	r := <-u.sendKeyPressCh
	num, err := strconv.Atoi(string(r))
	if err != nil {
		return u.pressAnyNumber()
	}
	return num
}
