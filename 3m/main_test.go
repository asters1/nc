package main

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	a := "G0 X23.1 Y76.9 F300"
	b := "N100 X45.3 Y78.0 S80 M3"
	c := "G1 X67 Y8"
	d := "g nj m"
	fmt.Println(formatXY(a))
	fmt.Println(formatXY(b))
	fmt.Println(formatXY(c))
	fmt.Println(formatXY(d))
}
