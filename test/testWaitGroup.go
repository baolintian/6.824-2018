package main

import "sync"
import "fmt"
import "time"

func main(){
	
	cnt := 0
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i<10; i++{
		go func(){
			cnt++
			time.Sleep(5*time.Second)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(cnt)
}