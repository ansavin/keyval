package main

import (
	"fmt"
)

var storage map[string]string

/* interface keyValueStorage interface {
	set() string
	get() string
	del() string
}
*/

func main() {
	storage = make(map[string]string)
	storage["1"] = "data"
	fmt.Println(storage["1"])
}