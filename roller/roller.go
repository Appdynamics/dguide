package roller

import (
	"fmt"
	"io"
	"sync"
	"syscall"
	"time"
)

const (
	// 150ms per frame
	DEFAULT_FRAME_RATE = time.Millisecond * 150
)

var DefaultCharset = []string{"|", "/", "-", "\\"}

// ANSI escape codes for colors and reset
const (
	ResetColor = "\033[0m"
)

var Colors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Purple
	"\033[36m", // Cyan
	"\033[37m", // White
}

type Roller struct {
	sync.Mutex
	Title     string
	Charset   []string
	FrameRate time.Duration
	runChan   chan struct{}
	stopOnce  sync.Once
	Output    io.Writer
	NoTty     bool
}

// create roller object
func NewRoller(title string) *Roller {
	sp := &Roller{
		Title:     title,
		Charset:   DefaultCharset,
		FrameRate: DEFAULT_FRAME_RATE,
		runChan:   make(chan struct{}),
	}
	// if !terminal.IsTerminal(syscall.Stdout) {
	// 	sp.NoTty = true
	// }
	if !IsTerminal(syscall.Stdout) {
		sp.NoTty = true
	}
	return sp
}

// start a new roller, title can be an empty string
func StartNew(title string) *Roller {
	return NewRoller(title).Start()
}

// start roller
func (sp *Roller) Start() *Roller {
	go sp.writer()
	return sp
}

// set custom roller frame rate
func (sp *Roller) SetSpeed(rate time.Duration) *Roller {
	sp.Lock()
	sp.FrameRate = rate
	sp.Unlock()
	return sp
}

// set custom roller character set
func (sp *Roller) SetCharset(chars []string) *Roller {
	sp.Lock()
	sp.Charset = chars
	sp.Unlock()
	return sp
}

// stop and clear the roller
func (sp *Roller) Stop() {
	//prevent multiple calls
	sp.stopOnce.Do(func() {
		close(sp.runChan)
		sp.clearLine()
	})
}

// roller animation
func (sp *Roller) animate() {
	var out string
	for i := 0; i < len(sp.Charset); i++ {
		//color := Colors[i%len(Colors)]
		color := Colors[3]
		out = color + sp.Charset[i] + ResetColor + " " + sp.Title
		switch {
		case sp.Output != nil:
			fmt.Fprint(sp.Output, out)
			//fmt.Fprint(sp.Output, "\r"+out)
		case !sp.NoTty:
			fmt.Print(out)
			//fmt.Print("\r" + out)
		}
		time.Sleep(sp.FrameRate)
		sp.clearLine()
	}
}

// write out roller animation until runChan is closed
func (sp *Roller) writer() {
	sp.animate()
	for {
		select {
		case <-sp.runChan:
			return
		default:
			sp.animate()
		}
	}
}

// workaround for Mac OS < 10 compatibility
func (sp *Roller) clearLine() {
	if !sp.NoTty {
		fmt.Printf("\033[2K")
		fmt.Println()
		fmt.Printf("\033[1A")
	}
}
