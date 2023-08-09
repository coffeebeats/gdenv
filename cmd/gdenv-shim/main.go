package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	a := os.Args

	log.SetFlags(0)
	log.Printf("godot %s", strings.Join(a[1:], " "))
}
