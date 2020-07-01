package main

import (
	"fmt"
	"github.com/ASalimov/bar"
	"github.com/ttacon/chalk"
	"time"
)

func main() {
	n := 20
	b := bar.NewWithOpts(bar.WithDimensions(100, 100),
		bar.WithFormat(
			fmt.Sprintf(
				" %sbuilding...%s :percent :bar %s:eta %s ",
				chalk.Blue,
				chalk.Reset,
				chalk.Green,
				chalk.Reset)))

	fmt.Println()
	fmt.Println()

	for i := 1; i < n; i++ {
		b.Tick()
		l := i
		if l > 5 {
			l = 5
		}
		b.Interruptf("%d is even!", i)
		time.Sleep(100 * time.Millisecond)
	}

	b.Done()

	fmt.Println()
	fmt.Println()
}
