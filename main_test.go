package main

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	a := CheckName("./input/a.nc")
	b := CheckName("./input/a.NC")
	c := CheckName("./input/a-b-c.NC")
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
}
