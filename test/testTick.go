package main

import (
	"fmt"
	"time"
	"math/rand"
)

func statusUpdate() string { return "" }

func randDuration(min time.Duration) time.Duration{
	rand.Seed(time.Now().Unix())
	return time.Duration(rand.Int63())%min+min
}

func main() {
	//tick还是按照固定的是将间隔触发管道传输的事件
	c := time.Tick(randDuration(time.Second))
	start := time.Now()
	for next := range c {
		fmt.Printf("%v %s\n", next, statusUpdate())
		fmt.Println(time.Since(start))
		start = time.Now()
	}
}