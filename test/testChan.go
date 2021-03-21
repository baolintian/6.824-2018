package main

import "fmt"
import "time"

func main(){
	// chan一定要通过make进行声明
	var message = make(chan string)
	go func(){
		go func(){
			//time.Sleep(time.Second*3)
			for value := range message{
				fmt.Println(value)
			}
		}()
		time.Sleep(time.Second)
		message <- "tianbaolin"
		time.Sleep(time.Second)
		message <- "babydragon"
	}()
	time.Sleep(time.Second*10)
}