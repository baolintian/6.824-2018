package main

import "fmt"
import "encoding/json"
import "os"

type KV struct{
	key string
	value []string
}

func main(){

	
	test := make(map[string][]string)
	
	test["name"] = append(test["name"], "baolintian")
	test["name"] = append(test["name"], "babydragon")

	// test["name"] = "baolintian"
	// test["sex"] = "male"
	res, err := os.Create("test.json")
	if err != nil{
		fmt.Println("failure")
		return
	}
	enc := json.NewEncoder(res)
	enc.Encode(test)

	//res.Close()
	res1, _ := os.Open("test.json")
	dec := json.NewDecoder(res1)
	
	for{
		var x map[string][]string
		err := dec.Decode(&x)
		if err != nil{
			break
		}
		for key, value := range x{
			fmt.Println(key, value)
		}
	}
}