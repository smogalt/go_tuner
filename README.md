# go_tuner

simple internet radio client inspired by pyradio (https://github.com/coderholic/pyradio). written in go mostly as a learning project but it is still very usable. 

### keybinds
p - pause
m - mute
-/+ - vol up/down
j/down arrow - down
k/up arrow - up

## config
stations go in a csv file (~/.config/go_tuner/stations) no spaces. *default file provided*

name,url

## dependencies
- socat
- mpv

## install
you can 'go run main.go' or build it to a binary (go build main.go). you can also use a precompiled binary copy the default stations list to .config/go_tuner/stations

