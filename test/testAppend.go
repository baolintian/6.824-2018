package main

import "fmt"

func main(){

	x := []int{1, 2, 3, 4, 5}
	y := append([]int{}, x[1:4]...)
	fmt.Println(y)
	

}