package main

import (
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(4)
	now := time.Now()
	total := 10000000
	ch := make(chan bool, 1024)
	go func() {
		for i := 0; i < total; i++ {
			ch <- true
		}
		close(ch)
	}()

	c := false
	for a := range ch {
		c = c || a
	}
	println(int(float64(total)/time.Now().Sub(now).Seconds()), "rqs")
}
