package main

import "time"
import "math/rand"
import "fmt"

func randTime(modtime time.Duration) time.Duration{
	rand.Seed(time.Now().Unix())
	extra := time.Duration(rand.Int63())%modtime
	return modtime + extra
}

const electionTimeout = 1000 * time.Millisecond


func main(){
	timerCount := randTime(electionTimeout)
	var test *time.Timer
	test = time.NewTimer(timerCount)

	start := time.Now()
	for{
		select{
		case <- test.C:
				fmt.Println(time.Since(start))
				start = time.Now()
				test = time.NewTimer(randTime(electionTimeout))
		default:

		}
	}
}