package main


import "fmt"

func main(){
	var x []string
	x = append(x, "tianbaolin")
	x = append(x, "babydrgon")
	x = append(x, "minghaozhao")
	for _, s := range x{
		fmt.Println(s)
	}
}