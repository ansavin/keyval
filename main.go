package main

import (
	"bufio"
    "fmt"
	"os"
	"strings"
)

var storage map[string]string

/* interface keyValueStorage interface {
	set() string
	get() string
	del() string
}
*/
func upload_to_storage(data []string) {
	for index := 0; index < len(data); index++ {
		key := index
		//fmt.Println(index)
		if index++; index >= len(data) {
			//fmt.Println("foo")
			//fmt.Println(index)
			storage[data[key]] = ""
		} else {
			//fmt.Println("bar")
			//fmt.Println(index)
			storage[data[key]] = data[index]
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	storage = make(map[string]string)
	input, _ := reader.ReadString('\n')
	data := strings.Fields(input)
	//fmt.Println("len", len(data))
	upload_to_storage(data)
	for key, value := range storage {
		fmt.Printf("%v: %v\n", key, value)
	}
}