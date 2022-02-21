package main

import (
	"log"
	"math/rand"
	"time"

	"os"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

const trailSize = 10

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	width, height := getTermSize()
	m := matrix{
		screen: s,
		height: height,
		width:  width,
	}
	enqueueTimeEvent(m.screen)
	loop(m)
}

func enqueueTimeEvent(s tcell.Screen) {
	time.AfterFunc(time.Millisecond*200, func() {
		// requeue time event
		err := s.PostEvent(&tcell.EventTime{})
		if err != nil {
			log.Println(err)
			quit(s)
		}
	})
}

func quit(s tcell.Screen) {
	s.Fini()
	os.Exit(0)
}

func loop(m matrix) {
	// startY shoul:d be in the upper half of the screen
	startX, startY := rand.Intn(m.width), rand.Intn(m.height/2)
	for {
		// Update screen
		m.screen.Show()

		// Poll event
		ev := m.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventTime:
			next := m.fillWithChars(startX, startY, trailSize)
			if next {
				startY++
			} else {
				startX, startY = rand.Intn(m.width), rand.Intn(m.height/2)
			}
			m.screen.Sync()
			enqueueTimeEvent(m.screen)
		case *tcell.EventResize:
			m.width, m.height = ev.Size()
			m.screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit(m.screen)
			}
		}

	}
}

var (
	characters = []rune{
		'a', 'b', 'c', 'd', 'e', 'f',
	}
)

type matrix struct {
	screen tcell.Screen
	height int
	width  int
}

func (m matrix) fillWithChars(x, y, trail int) bool {
	if y >= (m.height + trail) {
		return false
	}
	// clears previous fields
	m.screen.SetContent(x, y-trail, ' ', nil, tcell.StyleDefault)
	m.screen.SetContent(x, y, randomChar(), nil, tcell.StyleDefault)
	return true
}

func randomChar() rune {
	random := rand.Intn(len(characters))
	return characters[random]
}

func getTermSize() (width, height int) {
	width, height, err := term.GetSize(0)
	if err != nil {
		log.Fatal(err)
	}
	return width, height
}
