package main

import (
	"bufio"
    "fmt"
	"os"
	"strings"
)

var storage map[string]string

func upload_to_storage(data []string) {
	for index := 0; index < len(data); index++ {
		key := index
		if index++; index >= len(data) {
			storage[data[key]] = ""
		} else {
			storage[data[key]] = data[index]
		}
	}
}

func inspect_storage() {
	for key, value := range storage {
		fmt.Printf("%v: %v\n", key, value)
	}
}

func find_data_by_key(data []string) (string, string) {
	if  len(data) <= 1 {
		fmt.Println("no storage key specified")
		return "", ""
	} else {
		key := data[1]
		value, ok := storage[key]
		if !ok {
			fmt.Printf("wrong storage key: %v \n", key)
			return "", ""
		}
		return key, value
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	storage = make(map[string]string)
	loop: 
		for {
			input, _ := reader.ReadString('\n')
			data := strings.Fields(input)
			switch cmd := data[0]; cmd {
				case "set":
					upload_to_storage(data[1:])
				case "show":
					inspect_storage()
				case "quit":
					break loop
				case "get":
					key, value := find_data_by_key(data)
					if key != "" && value != "" {
						fmt.Printf("%v: %v\n", key, value)
					}
				case "del":
					key, value := find_data_by_key(data)
					if key != "" && value != "" {
						delete(storage,key)
					}
				default:
					fmt.Printf("wrong cmd: %v \n", cmd)
					break loop
			}
		}
}