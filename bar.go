package bar

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

var noop = func() {}

// Bar is a progress bar to be used for displaying task progress
// via terminal output
type Bar struct {
	progress, total, width     int
	start, end                 string
	complete, head, incomplete string
	closed                     bool
	startedAt                  time.Time
	rate                       float64
	eta                        time.Duration
	formatString               string
	format                     []token
	context                    []*ContextValue
	callback                   func()
	output                     Output
	debug                      bool
	buffer                     []string
	lines                      int
	curPosition                int
}

// ContextValue is a tuple that defines a substitution for a custom verb
type ContextValue struct {
	verb  string
	value *stringish
}

// Context is a wrapper type for a slice of ContextValues
type Context []*ContextValue

// Ctx is a helper for creating a ContextValue tuple
func Ctx(verb string, value interface{}) *ContextValue {
	if verb[0] == ':' {
		panic(fmt.Sprintf("don't prefix your custom verb declaration with a `:`, it's implied (at %s)", verb))
	}

	if verb == "bar" || verb == "percent" || verb == "rate" || verb == "eta" {
		panic(fmt.Sprintf(":%s is a reserved verb, please choose another name", verb))
	}

	return &ContextValue{
		verb:  verb,
		value: newStringish(value),
	}
}

const bufferSize = 20
const defaultFormat = " :bar :percent :rate ops/s "

// New creates a new instance of bar.Bar with the given total and
// returns a reference to it
func New(t int) *Bar {
	return &Bar{
		progress:     0,
		total:        t,
		width:        20,
		start:        "(",
		complete:     "█",
		head:         "█",
		incomplete:   " ",
		end:          ")",
		closed:       false,
		startedAt:    time.Now(),
		rate:         0,
		formatString: defaultFormat,
		format:       tokenize(defaultFormat, []string{}),
		callback:     noop,
		output:       &stdout{},
		buffer:       []string{},
		lines:        0,
	}
}

// Tick increments the bar's progress by 1
func (b *Bar) Tick() {
	if !b.canUpdate("Tick") {
		return
	}

	b.TickAndUpdate(nil)
}

// TickAndUpdate is a helper function for calling Tick
// followed by Update
func (b *Bar) TickAndUpdate(ctx Context) {
	if !b.canUpdate("TickAndUpdate") {
		return
	}

	b.Update(b.progress+1, ctx)
}

// Update sets the bar's progress to an arbitrary value
// and optionally updates the bar's context
func (b *Bar) Update(progress int, ctx Context) {
	if !b.canUpdate("Update") {
		return
	}

	duration := time.Now().Sub(b.startedAt)
	b.rate = float64(b.progress) / duration.Seconds()
	b.eta = time.Duration(float64(b.total-b.progress)/b.rate) * time.Second

	b.progress = progress

	if ctx != nil {
		b.context = ctx
		b.format = tokenize(b.formatString, ctx.customVerbs())
	}

	b.write()
}

func (b *Bar) SetFormat(f string) {
	b.format = tokenize(f, nil)
}

// Done finalizes the bar and prints it followed by a new line
func (b *Bar) Done() {
	b.closed = true
	b.write()
	fmt.Println()
	b.callback()
}

// SetLines
func (b *Bar) SetLines(lines int) {

	for i := 0; i < lines; i++ {
		fmt.Print("\033[F")
	}
	if b.curPosition > 0 {
		for i := b.curPosition; i >= 0; i-- {
			b.output.ClearLine()
			for {
				if len(b.buffer) < i+1 {
					b.buffer = append(b.buffer, "")
				} else {
					break
				}
			}

			fmt.Println(b.buffer[i])
		}
	}

	b.curPosition = b.curPosition + (lines - b.lines)
	//b.output.ClearLine()\
	b.lines = lines
	b.write()
}

func (b *Bar) GetLines() int {
	return b.lines
}

// Interruptf passes the given input to fmt.Sprintf and prints
// it above the bar
func (b *Bar) Interruptf(format string, s ...interface{}) {
	b.Interrupt(fmt.Sprintf(format, s...))
}

// Interrupt prints s above the bar
func (b *Bar) Interrupt(s string) {
	if b.closed {
		return
	}
	if b.lines == 0 {
		b.output.ClearLine()
		fmt.Println(s)
		b.write()
		return
	}
	if b.curPosition == b.lines {
		for i := 0; i < b.lines; i++ {
			fmt.Print("\033[F")
		}
		if len(b.buffer) < bufferSize {
			b.buffer = append(b.buffer, b.buffer[len(b.buffer)-1])
		}
		for i := len(b.buffer) - 2; i >= 0; i-- {
			b.buffer[i+1] = b.buffer[i]
		}
		for i := b.lines - 1; i > 0; i-- {
			b.output.ClearLine()
			fmt.Println(b.buffer[i])
		}
	} else {
		b.curPosition++
		b.buffer = append(b.buffer, "")
		for i := len(b.buffer) - 2; i >= 0; i-- {
			b.buffer[i+1] = b.buffer[i]
		}
	}
	b.buffer[0] = s
	b.output.ClearLine()
	fmt.Println(s)
	b.write()
}

func (b *Bar) write() {
	b.output.ClearLine()
	b.output.Printf("%s", b)
}

func (b *Bar) canUpdate(method string) bool {
	if b.closed {
		fmt.Fprintf(os.Stderr, "bar: attempted to call %s on a closed bar, this is likely caused by a memory leak", method)
		return false
	}

	return true
}

func (b *Bar) prog() float64 {
	return float64(b.progress) / float64(b.total)
}

func (c Context) customVerbs() []string {
	verbs := make([]string, len(c))

	for _, def := range c {
		verbs = append(verbs, def.verb)
	}

	return verbs
}

func (b *Bar) String() string {
	var buf bytes.Buffer

	for _, s := range b.format {
		if b.debug {
			buf.WriteString(s.debug(b))
		} else {
			buf.WriteString(s.print(b))
		}
	}

	return buf.String()
}
