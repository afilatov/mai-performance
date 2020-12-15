package main
// ulimit -c unlimited
// GOTRACEBACK=crash ./coredumper
// dlv core ./coredumper ./core

import (
	"fmt"
	"math/rand"
)

func main() {
	var sum int
	for {
		n := rand.Intn(1e6)
		sum += n
		if sum % 42 == 0 {
			fmt.Println(sum)
			panic("bad sum")
		}
	}
}
