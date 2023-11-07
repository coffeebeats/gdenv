package progress

import (
	"testing"
)

/* ------------------------------ Test: Writer ------------------------------ */

func TestWriter(t *testing.T) {
	total := uint64(64)

	// Given: A 'Progress' struct to report progress through.
	p, err := NewWithTotal(total)
	if err != nil {
		t.Errorf("err: got %#v, want %#v", err, nil)

	}

	// Given: Channels to communicate write signals through.
	write, wrote := make(chan struct{}), make(chan struct{})

	// Given: A goroutine that writes a single byte to a 'Writer' when signaled.
	go func(p *Progress) {
		// Given: A new 'Writer' reporting progress via the 'Progress' struct.
		w := NewWriter(p)

		for range write {
			// When: '1' byte is written to the reporter.
			n, err := w.Write([]byte{1})

			// Then: There's no error writing to the progress reporter.
			if err != nil {
				t.Errorf("err: got %#v, want %#v", err, nil)
			}

			// Then: The returned number of bytes written is correct.
			if n != 1 {
				t.Errorf("output: got %d, want %d", n, 1)
			}

			wrote <- struct{}{}
		}

		close(wrote)
	}(p)

	for i := range make([]struct{}, total) {
		// Given: The correct initial progress value.
		if got, want := p.Percentage(), float64(i)/float64(total); got != want {
			t.Errorf("output: got %#v, want %#v", got, want)
		}

		// When: A single byte is written in another thread.
		write <- struct{}{}

		<-wrote // Wait for the other thread to write.

		// Then: The 'Progress' value in this thread updates accordingly.
		if got, want := p.Percentage(), float64(i+1)/float64(total); got != want {
			t.Errorf("output: got %#v, want %#v", got, want)
		}
	}

	close(write)
}
