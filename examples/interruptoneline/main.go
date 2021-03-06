package main

import (
	"fmt"
	"github.com/ASalimov/bar"
	"github.com/ttacon/chalk"
	"os"
	"strings"
	"time"
)

func main() {
	bar.InitTerminal()
	ch := make(chan string)
	b := bar.NewWithOpts(bar.WithDimensions(100, 100),
		bar.WithLines(2),
		bar.WithFormat(
			fmt.Sprintf(
				" %sbuilding...%s :percent :bar %s:eta %s ",
				chalk.Blue,
				chalk.Reset,
				chalk.Green,
				chalk.Reset)))

	fmt.Println()
	fmt.Println()
	listen_keys(ch)
	i := 1
	for {
		//b.Tick()
		l := i
		if l > 5 {
			l = 5
		}
		select {
		case stdin, _ := <-ch:
			//fmt.Println("Keys pressed:", []byte(stdin))
			if []byte(stdin)[0] == 10 {
				b.SetLines(b.GetLines() + 1)
				//time.Sleep(3000 * time.Millisecond)
			} else {
			}
			//time.Sleep(1000 * time.Millisecond)
			//time.Sleep(5000 * time.Millisecond)
			//b.Interruptf("%d is even!", i)
		case <-time.After(200 * time.Millisecond):
			b.Interruptf("%d is even!"+strings.Repeat(".", int(100)-i), i)
			i++
		}

	}

	b.Done()

	fmt.Println()
	fmt.Println()
}

//
//
func listen_keys(ch chan string) {

	go func(ch chan string) {
		// disable input buffering
		//exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		//// do not display entered characters on the screen
		//exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch)

}
