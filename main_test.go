package main

import "testing"

func Test_greet(t *testing.T) {
	want := "Hello, world!"
	if got := greet("world"); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
