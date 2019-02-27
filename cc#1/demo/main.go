package main

import (
	"fmt"
	"mockdemo"
)

func main() {
	counter := mockdemo.NewPageCounter()
	if n, err := counter.Count("https://www.process-one.net", "ProcessOne"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Count:", n)
	}
}
