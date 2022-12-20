# Vmgo

A virtual machine implementation following the tutorial of [Write your Own Virtual Machine](https://www.jmeiners.com/lc3-vm/), but written in Go instead of C.

## Usage
In the original tutorial there are 2 programs available that can be run in the virtual machine: [2048.obj](https://www.jmeiners.com/lc3-vm/supplies/2048.obj) and [rogue.obj](https://www.jmeiners.com/lc3-vm/supplies/rogue.obj).

Both are also here and you can:

## Build
```bash
go build main.go
./main -image 2048.obj
./main -image rogue.obj
```

## Run
```bash
go run main.go -image 2048.obj
go run main.go -image rogue.obj
```

## Notes
Some concepts I couldn't reproduce in a similar way in Go. In these cases, I consulted examples and learned a few things, such as:

- Keyboard

 Since it's not relevant to virtual machines, the tutorial gives you the code to access the keyboard for copy and paste. I first tried to implement similar behavior using [bufio](https://pkg.go.dev/bufio), but switched to the [keyboard library](https://github.com/eiannone/keyboard) afterwards.
