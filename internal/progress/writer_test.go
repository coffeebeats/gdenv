package progress

import (
	"errors"
	"fmt"
	"testing"
)

/* ------------------------------ Test: Writer ------------------------------ */

func TestWriter(t *testing.T) {
	// Given: A 'Progress' struct to report progress through.
	p, err := New(4)
	if err != nil {
		t.Fatalf("err: got %#v, want %#v", err, nil)

	}

	// Given: Channels to communicate write signals through.
	write, wrote, errs := make(chan struct{}), make(chan struct{}), make(chan error)

	// Given: A goroutine that writes a single byte to a 'Writer' when signaled.
	go func(p *Progress) {
		// Given: A new 'Writer' reporting progress via the 'Progress' struct.
		w := NewWriter(p)

		for {
			select {
			case <-write:
				n, err := w.Write([]byte{1})
				if err != nil {
					errs <- err
				}
				if n != 1 {
					errs <- fmt.Errorf("output: got %#v, want %#v", n, 1)
				}

				wrote <- struct{}{}
			default: // channel 'write' closed
				close(wrote)
				return
			}
		}
	}(&p)

	for i := range [4]int{} {
		// Given: The correct initial progress value.
		got, err := p.Percentage()
		if !errors.Is(err, nil) {
			t.Fatalf("err: got %#v, want %#v", err, nil)
		}
		if want := float64(i) / float64(4); got != want {
			t.Fatalf("output: got %#v, want %#v", got, want)
		}

		// When: A single byte is written in another thread.
		write <- struct{}{}

		select {
		case <-wrote:
			got, err := p.Percentage()
			if !errors.Is(err, nil) {
				t.Fatalf("err: got %#v, want %#v", err, nil)
			}
			// Then: The 'Progress' value in this thread updates accordingly.
			if want := float64(i+1) / float64(4); got != want {
				t.Fatalf("output: got %#v, want %#v", got, want)
			}
		case err := <-errs:
			t.Fatalf("err: got %#v, want %#v", err, nil)
		}
	}

	close(write)
}