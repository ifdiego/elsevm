# Vmgo

A virtual machine implementation following the tutorial of [Write your Own Virtual Machine](https://www.jmeiners.com/lc3-vm/), but written in Go instead of C.

## Usage
The original tutorial has 2 programs available to run in the virtual machine: [2048.obj](https://www.jmeiners.com/lc3-vm/supplies/2048.obj) and [rogue.obj](https://www.jmeiners.com/lc3-vm/supplies/rogue.obj).

## Build
```bash
go build main.go
./main 2048.obj
./main rogue.obj
```

## Run
```bash
go run main.go 2048.obj
go run main.go rogue.obj
```

## Notes
The code to access the keyboard in the tutorial is made to copy and paste. I tried to implement a similar behavior, first using [bufio](https://pkg.go.dev/bufio) but switching to [keyboard library](https://github.com/eiannone/keyboard) afterwards.
