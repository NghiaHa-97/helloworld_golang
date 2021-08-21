package main

import (
	"encoding/json"
	"fmt"
)

type employee struct {
	Name string `json:" myName "`
	Age  int    `json:" age "`
}

func main() {
	a := make(map[string]employee)
	a["1"] = employee{Name: "John", Age: 24}
	j, err := json.Marshal(a)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Println(string(j))
	}
}
