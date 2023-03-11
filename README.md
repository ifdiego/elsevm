# elsevm

A virtual machine implementation written in Go.

The main purpose was learning what is, how does it work and what could be done
with a virtual machine by following the [Write your Own Virtual
Machine](https://www.jmeiners.com/lc3-vm/) tutorial without copy-pasting the
entire code, that's the reason why I used another language, Go instead of C.

#### Usage

Both [2048.obj](https://github.com/ifdiego/elsevm/blob/main/2048.obj) and
[rogue.obj](https://github.com/ifdiego/elsevm/blob/main/rogue.obj) assembly
programs are available.

In project's directory, compile:

```bash
go build
```

Run it, passing an assembly program as argument:

```bash
./elsevm 2048.obj
./elsevm rogue.obj
```

You can also run main file without compiling:

```bash
go run main.go 2048.obj
go run main.go rogue.obj
```

Feel free to test other assemblies.

#### Access the keyboard

Originally, `termios.h` was used in the tutorial to deal with
input/output communication ports.

I have tried a similar behavior using [bufio](https://pkg.go.dev/bufio), a Go's
standard library package. But I couldn't fix some bugs, which led me switching
to [eiannone keyboard library](https://github.com/eiannone/keyboard) afterwards.
