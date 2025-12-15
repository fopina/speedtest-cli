package oututil

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	"golang.org/x/sync/errgroup"
)

var outTTY = false

func init() {
	outTTY = isatty.IsTerminal(uintptr(os.Stdout.Fd()))
}

type Printer struct {
	c chan<- string
	g errgroup.Group
}

func StartPrinting() *Printer {
	if !outTTY {
		return nil
	}
	c := make(chan string)
	p := Printer{c: c}
	p.g.Go(func() error {
		s, ok := <-c
		if !ok {
			return nil
		}
		fmt.Println(s)
		for s := range c {
			erasePrevious()
			fmt.Println(s)
		}
		return nil
	})
	return &p
}

func (p *Printer) Println(s string) {
	if p == nil {
		return
	}
	p.c <- s
}

func (p *Printer) Finalize(s string) {
	if p == nil {
		fmt.Println(s)
		return
	}
	close(p.c)
	p.g.Wait()
	if s == "" {
		return
	}
	erasePrevious()
	fmt.Println(s)
}

func erasePrevious() {
	fmt.Print("\033[1A") // ANSI escape sequence for move one line up
	fmt.Print("\033[K")  // ANSI escape sequence for erase current line
}
