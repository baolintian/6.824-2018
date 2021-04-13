package main

import "fmt"

func main(){

	//var x map[string]bool
	x := make(map[string]bool)
	x["name"] = true
	fmt.Println(x["name"])
	fmt.Println(x["sex"])
	var unique []string
	unique = append(unique, "name")
	unique = append(unique, "sex")
	fmt.Println(unique)
	for index, k := range x{
		fmt.Println(index, k)
	}
}