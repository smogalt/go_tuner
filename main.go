package main

import (
	"encoding/csv"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"os"
	"os/exec"
)

type item struct {
	name string
	url  string
}

func main() {
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	home_dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	stations_raw, err := os.Open(home_dir + "/.config/go_tuner/stations")
	if err != nil {
		panic(err)
	}
	defer stations_raw.Close()
	csvReader := csv.NewReader(stations_raw)
	data, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	stations := parse_csv(data)

	selected := 0
	playing := -1

	err = termbox.Init()
	termbox.HideCursor()
	if err != nil {
		panic(err)
	}

	box()
	print_menu(selected, stations)
	termbox.Flush()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch {
				case ev.Key == termbox.KeyArrowDown:
					if selected == len(stations)-1 {
						selected = -1
					}
					selected++

				case ev.Key == termbox.KeyArrowUp:
					if selected == 0 {
						selected = len(stations)
					}
					selected--

				case ev.Ch == 'q':
					quit()

				case ev.Ch == 'p':
					pause()

				case ev.Key == termbox.KeyEnter:
					if playing != selected {
						play(stations[selected])
						playing = selected
					}

				}

			}
			refresh()
			print_menu(selected, stations)
			termbox.Flush()
		}

	}
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func box() {
	w, h := termbox.Size()
	for i := 1; i < w-1; i++ {
		termbox.SetCell(i, 0, '─', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(i, h-1, '─', termbox.ColorWhite, termbox.ColorDefault)
	}

	for j := 1; j < h-1; j++ {
		termbox.SetCell(1, j, '│', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(w-2, j, '│', termbox.ColorWhite, termbox.ColorDefault)
	}

	termbox.SetCell(1, 0, '╭', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(w-2, 0, '╮', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(w-2, h-1, '╯', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(1, h-1, '╰', termbox.ColorWhite, termbox.ColorDefault)

	tbprint(4, 0, termbox.ColorWhite, termbox.ColorDefault, " go tuner ")
}

func print_menu(sel int, list [16]item) {
	for i := 0; i < len(list); i++ {
		if i == sel {
			tbprint(4, (i*2)+2, termbox.ColorRed, termbox.ColorWhite, list[i].name)
		} else {
			tbprint(4, (i*2)+2, termbox.ColorRed, termbox.ColorDefault, list[i].name)
		}
	}
}

func refresh() {
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)
	box()
	termbox.Flush()
}

func quit() {
	cmd := exec.Command("/bin/sh", "-c", `echo '{"command": ["quit"]}' | socat - /tmp/mpvsocket`)
	cmd.Run()
	termbox.Close()
	os.Exit(0)
}

func play(s item) {
	cmd := exec.Command("/bin/sh", "-c", `echo '{"command": ["quit"]}' | socat - /tmp/mpvsocket`)
	cmd.Run()
	cmd = exec.Command("mpv", s.url, "--input-ipc-server=/tmp/mpvsocket")
	cmd.Start()
}

func pause() {
	cmd := exec.Command("/bin/sh", "-c", `echo '{"command": ["cycle", "pause"]}' | socat - /tmp/mpvsocket`)
	cmd.Run()
}

func parse_csv(data [][]string) [16]item {
	var menu_items [16]item
	for i, line := range data {
		if i < len(menu_items) {
			for j, field := range line {
				if j == 0 {
					menu_items[i].name = field
				} else if j == 1 {
					menu_items[i].url = field
				}
			}
		}
	}
	return menu_items
}
