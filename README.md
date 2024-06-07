# go_tuner

simple internet radio client inspired by pyradio (https://github.com/coderholic/pyradio). written in go mostly as a learning project but it is still very usable. 

### keybinds
enter/space - select

q - quit
<<<<<<< HEAD

up arrow - up on the list

down arrow - down on the list

enter - select station

=======
>>>>>>> 9c4ac1c (Added default station list. Updated README. Added two column display and various keybinds)
p - pause
m - mute
-/+ - vol up/down
j/down arrow - down
k/up arrow - up

## config

<<<<<<< HEAD
stations go in a csv file (~/.config/go_tuner/stations) no spaces
(*i'll add a default file eventually*)
=======
stations go in a csv file (~/.config/go_tuner/stations) no spaces. *default file provided*

name,url

## dependencies
- socat
- mpv

## install
you can 'go run main.go' or build it to a binary (go build main.go). you can also use a precompiled binary copy the default stations list to .config/go_tuner/stations

>>>>>>> 9c4ac1c (Added default station list. Updated README. Added two column display and various keybinds)
