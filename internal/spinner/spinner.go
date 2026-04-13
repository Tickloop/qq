// Package spinner provides terminal loading indicators.
package spinner

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Spinner controls a terminal loading animation.
// Start begins the animation; Stop halts it and cleans up.
// Stop is safe to call multiple times or without a preceding Start.
type Spinner interface {
	Start()
	Stop()
}

// braille dot frames for the animation cycle
var frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// ANSISpinner renders a green braille-dot animation to an io.Writer.
type ANSISpinner struct {
	w        io.Writer
	interval time.Duration
	mu       sync.Mutex
	done     chan struct{}
	stopped  bool
	wg       sync.WaitGroup
}

// NewANSISpinner creates a spinner that writes ANSI-colored frames to w
// at the given interval. Call Start to begin animation.
func NewANSISpinner(w io.Writer, interval time.Duration) *ANSISpinner {
	return &ANSISpinner{
		w:        w,
		interval: interval,
	}
}

// Start begins the spinner animation in a background goroutine.
func (s *ANSISpinner) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.done = make(chan struct{})
	s.stopped = false

	done := s.done
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.run(done)
	}()
}

func (s *ANSISpinner) run(done <-chan struct{}) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// hide cursor while spinner is active
	fmt.Fprint(s.w, "\033[?25l")

	i := 0
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			fmt.Fprintf(s.w, "\r\033[32m%s\033[0m ", frames[i%len(frames)])
			i++
		}
	}
}

// Stop halts the animation and restores the terminal line.
// Safe to call multiple times or without a preceding Start.
func (s *ANSISpinner) Stop() {
	s.mu.Lock()
	if s.stopped || s.done == nil {
		s.mu.Unlock()
		return
	}
	s.stopped = true
	close(s.done)
	s.mu.Unlock()

	// wait for the goroutine to finish before clearing the line
	s.wg.Wait()

	// clear spinner line and restore cursor
	fmt.Fprint(s.w, "\r\033[K\033[?25h")
}

// NopSpinner is a no-op Spinner for non-TTY environments or testing.
type NopSpinner struct{}

func (NopSpinner) Start() {}
func (NopSpinner) Stop()  {}
