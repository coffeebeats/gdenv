package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	a := os.Args
	fmt.Println("godot", strings.Join(a[1:], " "))
}
