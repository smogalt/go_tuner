package main

import (
	"encoding/csv"
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"os"
	"os/exec"
)

type station struct {
	name string
	url  string
}

type playlist struct {
	stations     []station
	selected     int
	playing      int
	volume       int
	volume_human string
}

func main() {
	/* playlist struct to hold the stations and which ones are playing */
	var pl playlist

	/* event queue for keyboard input */
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	/* init default station and load stations from ~/.config/go_tuner/stations */
	pl.stations = load_stations()
	pl.selected = 0
	pl.playing = -1
	pl.volume = 100
	pl.volume_human = fmt.Sprint(100)

	/* start termbox */
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	/* set screen */
	refresh_screen(pl)

	/* input loop */
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch {
				/* down in list */
				case ev.Key == termbox.KeyArrowDown || ev.Ch == 'j':
					if pl.selected == len(pl.stations)-1 {
						pl.selected = 0
					} else {
						pl.selected++
					}
				/* up in list */
				case ev.Key == termbox.KeyArrowUp || ev.Ch == 'k':
					if pl.selected == 0 {
						pl.selected = len(pl.stations) - 1
					} else {
						pl.selected--
					}

				/* quit */
				case ev.Ch == 'q':
					quit()

				/* pause */
				case ev.Ch == 'p':
					pause()

				/* select station */
				case ev.Key == termbox.KeyEnter || ev.Key == termbox.KeySpace:
					play(&pl)
					if pl.selected == len(pl.stations)-1 {
						pl.selected = 0
					} else {
						pl.selected++
					}

				case ev.Ch == '-' || ev.Ch == '_':
					vol(-1, &pl)

				case ev.Ch == '=' || ev.Ch == '+':
					vol(1, &pl)

				case ev.Ch == 'm':
					vol(0, &pl)
				}
			}
			refresh_screen(pl)
		}
	}
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func print_box(p playlist) {
	w, h := termbox.Size()
	/* horizontal sides */
	for i := 1; i < w-1; i++ {
		termbox.SetCell(i, 0, '─', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(i, h-1, '─', termbox.ColorWhite, termbox.ColorDefault)
	}

	/* vertical sides */
	for j := 1; j < h-1; j++ {
		termbox.SetCell(1, j, '│', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(w-2, j, '│', termbox.ColorWhite, termbox.ColorDefault)
	}

	/* corners */
	termbox.SetCell(1, 0, '╭', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(w-2, 0, '╮', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(w-2, h-1, '╯', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(1, h-1, '╰', termbox.ColorWhite, termbox.ColorDefault)

	/* title */
	tbprint(4, 0, termbox.ColorWhite, termbox.ColorDefault, " go tuner "+"─"+" vol: "+p.volume_human+" ")
}

func print_menu(p playlist) {
	x := 5
	y := 2
	w, h := termbox.Size()
	for i := 0; i < len(p.stations); i++ {
		y = ((i * 2) + 2) % (h - 4)

		if ((i * 2) + 2) >= (h - 4) {
			x = (w / 2) + 5
			if h%2 == 0 {
				y += 2
			} else {
				y++
			}
		}

		if i == p.playing {
			tbprint(x, y, termbox.ColorLightRed, termbox.ColorDefault, p.stations[i].name)
			tbprint(x-2, y, termbox.ColorDarkGray, termbox.ColorDefault, "*")
		} else {
			tbprint(x, y, termbox.ColorWhite, termbox.ColorDefault, p.stations[i].name)
		}
		if i == p.selected {
			tbprint(x, y, termbox.ColorCyan, termbox.ColorDefault, p.stations[i].name)
			tbprint(x-2, y, termbox.ColorLightCyan, termbox.ColorDefault, ">")
		}
		if i == p.selected && p.selected == p.playing {
			tbprint(x, y, termbox.ColorCyan, termbox.ColorDefault, p.stations[i].name)
		}
	}
}

func print_help() {
	_, h := termbox.Size()
	tbprint(4, h-3, termbox.ColorDarkGray, termbox.ColorDefault, "k - up | j - down | enter/space - select")
	tbprint(4, h-2, termbox.ColorDarkGray, termbox.ColorDefault, "m - mute | p - pause | q - quit | -/+ - vol up/down ")
}

func refresh_screen(p playlist) {
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)
	print_box(p)
	print_menu(p)
	print_help()
	termbox.Flush()
	w, h := termbox.Size()
	termbox.SetCursor(w, h)
}

func quit() {
	cmd := exec.Command("/bin/sh", "-c", `echo '{"command": ["quit"]}' | socat - /tmp/mpvsocket`)
	cmd.Run()
	termbox.Close()
	os.Exit(0)
}

func play(p *playlist) {
	if p.selected != p.playing {
		cmd := exec.Command("/bin/sh", "-c", `echo '{"command": ["quit"]}' | socat - /tmp/mpvsocket`)
		cmd.Run()
		cmd = exec.Command("mpv", p.stations[p.selected].url, "--input-ipc-server=/tmp/mpvsocket")
		cmd.Start()
		p.playing = p.selected
	} else {
		pause()
		return
	}
}

func pause() {
	cmd := exec.Command("/bin/sh", "-c", `echo '{"command": ["cycle", "pause"]}' | socat - /tmp/mpvsocket`)
	cmd.Run()
}

/* true is up false is down */
func vol(i int, p *playlist) {
	if i == 1 && p.volume <= 95 {
		p.volume += 5
	}
	if i == -1 && p.volume >= 5 {
		p.volume -= 5
	}
	if i == 0 {
		if p.volume_human == "muted" {
			p.volume_human = fmt.Sprint(p.volume)
		} else {
			p.volume_human = "muted"
		}

		cmd := exec.Command("/bin/sh", "-c", `echo '{"command": ["cycle", "mute"]}' | socat - /tmp/mpvsocket`)
		cmd.Run()
		return
	}
	cmd := exec.Command("/bin/sh", "-c", `echo '{"command": ["set_property", "volume", `+fmt.Sprint(p.volume)+`]}' | socat - /tmp/mpvsocket`)
	cmd.Run()

	if p.volume == 0 {
		p.volume_human = "muted"
	} else {
		p.volume_human = fmt.Sprint(p.volume)
	}
}

func load_stations() []station {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	stations_raw, err := os.Open(home_dir + "/.config/go_tuner/stations")
	if err != nil {
		cmd := exec.Command("mkdir", "-p", home_dir+"/.config/go_tuner/")
		cmd.Run()
		cmd = exec.Command("touch", home_dir+"/.config/go_tuner/stations")
		cmd.Run()
	}

	csvReader := csv.NewReader(stations_raw)
	data, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	var menu_items []station

	for i, line := range data {
		if i > 0 {
			var temp station
			for j, field := range line {
				if j == 0 {
					temp.name = field
				} else if j == 1 {
					temp.url = field
				}
			}
			menu_items = append(menu_items, temp)
		}
	}

	stations_raw.Close()

	return menu_items
}
