package main

import (
	"fmt"
	"time"
)

func main() {
	a := time.Now()
	b := a

	fmt.Println(a == b)
}
