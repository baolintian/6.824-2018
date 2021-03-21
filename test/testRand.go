package main

import "fmt"
import "math/rand"
import "time"

func main(){
	
	rand.Seed(time.Now().Unix())
	fmt.Println(rand.Intn(10))
}