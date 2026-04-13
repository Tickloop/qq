package spinner

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestANSISpinner_SatisfiesInterface(t *testing.T) {
	var _ Spinner = &ANSISpinner{}
}

func TestNopSpinner_SatisfiesInterface(t *testing.T) {
	var _ Spinner = NopSpinner{}
}

func TestANSISpinner_RendersGreenFramesAndCleanup(t *testing.T) {
	var buf bytes.Buffer
	s := NewANSISpinner(&buf, 10*time.Millisecond)

	s.Start()
	time.Sleep(200 * time.Millisecond)
	s.Stop()

	out := buf.String()

	// should contain hide-cursor sequence at the start
	if !strings.Contains(out, "\033[?25l") {
		t.Error("missing hide-cursor sequence")
	}

	// should contain at least one green frame
	if !strings.Contains(out, "\033[32m") {
		t.Error("missing green color code")
	}
	if !strings.Contains(out, "\033[0m") {
		t.Error("missing color reset code")
	}

	// should end with cleanup: clear line + show cursor
	if !strings.HasSuffix(out, "\r\033[K\033[?25h") {
		t.Error("missing cleanup sequence at end of output")
	}
}

func TestANSISpinner_StopIsIdempotent(t *testing.T) {
	var buf bytes.Buffer
	s := NewANSISpinner(&buf, 10*time.Millisecond)

	s.Start()
	time.Sleep(50 * time.Millisecond)

	// calling Stop multiple times must not panic
	s.Stop()
	s.Stop()
	s.Stop()
}

func TestANSISpinner_StopWithoutStartIsNoOp(t *testing.T) {
	var buf bytes.Buffer
	s := NewANSISpinner(&buf, 10*time.Millisecond)

	// Stop without Start must not panic or write anything
	s.Stop()

	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestANSISpinner_CleanupAppearsOnce(t *testing.T) {
	var buf bytes.Buffer
	s := NewANSISpinner(&buf, 10*time.Millisecond)

	s.Start()
	time.Sleep(100 * time.Millisecond)
	s.Stop()

	out := buf.String()
	showCursor := "\033[?25h"
	count := strings.Count(out, showCursor)
	if count != 1 {
		t.Errorf("expected show-cursor sequence exactly once, got %d", count)
	}
}

func TestANSISpinner_FrameCharactersAreValid(t *testing.T) {
	var buf bytes.Buffer
	s := NewANSISpinner(&buf, 10*time.Millisecond)

	s.Start()
	time.Sleep(200 * time.Millisecond)
	s.Stop()

	out := buf.String()
	valid := map[string]bool{}
	for _, f := range frames {
		valid[f] = true
	}

	// extract characters between green start and reset sequences
	parts := strings.Split(out, "\033[32m")
	count := 0
	for _, p := range parts[1:] { // skip text before first green code
		ch, _, ok := strings.Cut(p, "\033[0m")
		if !ok {
			continue
		}
		if !valid[ch] {
			t.Errorf("unexpected frame character: %q", ch)
		}
		count++
	}
	if count == 0 {
		t.Error("no frame characters found in output")
	}
}

func TestANSISpinner_ConcurrentStartStop(t *testing.T) {
	// each goroutine gets its own spinner to verify internal race-freedom
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var buf bytes.Buffer
			s := NewANSISpinner(&buf, 10*time.Millisecond)
			s.Start()
			time.Sleep(20 * time.Millisecond)
			s.Stop()
		}()
	}
	wg.Wait()
}

func TestNopSpinner_StartStopAreNoOps(t *testing.T) {
	s := NopSpinner{}
	// must not panic
	s.Start()
	s.Stop()
	s.Start()
	s.Stop()
}
