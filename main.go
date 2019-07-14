package main

import (
	"bufio"
    "fmt"
	"net"
	"strings"
)

var storage map[string]string

var addr string = ":1234"

func handle(conn net.Conn) error {
    defer func() {
        fmt.Printf("Closing connection from %v \n", conn.RemoteAddr())
        conn.Close()
    }()
    r := bufio.NewReader(conn)
    w := bufio.NewWriter(conn)
    scanr := bufio.NewScanner(r)
	loop:
		for {
			scanned := scanr.Scan()
			if !scanned {
				if err := scanr.Err(); err != nil {
					fmt.Printf("%v(%v)\n", err, conn.RemoteAddr())
					return err
				}
				break
			}
			data := strings.Fields(scanr.Text())
			switch cmd := data[0]; cmd {
				case "set":
					upload_to_storage(data[1:])
				case "show":
					for key, value := range storage {
						w.WriteString(fmt.Sprintf("%v: %v\n", key, value))
						}
					w.Flush()
				case "quit":
					break loop
				case "get":
					key, value := find_data_by_key(data)
					if key != "" && value != "" {
						w.WriteString(fmt.Sprintf("%v: %v\n", key, value))
						w.Flush()
					} else {
						w.WriteString(fmt.Sprintf("storage key not found or not specified\n"))
						w.Flush()
					}
				case "del":
					key, value := find_data_by_key(data)
					if key != "" && value != "" {
						delete(storage,key)
					} else {
						w.WriteString(fmt.Sprintf("storage key not found or not specified\n"))
						w.Flush()
					}
				default:
					w.WriteString(fmt.Sprintf("wrong cmd: %v \n", cmd))
					w.Flush()
					break loop
			}
		}
    return nil
}

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

func find_data_by_key(data []string) (string, string) {
	if  len(data) <= 1 {
		fmt.Println("Server error: no storage key specified\n")
		return "", ""
	} else {
		key := data[1]
		value, ok := storage[key]
		if !ok {
			fmt.Printf("Server error: wrong storage key: %v \n", key)
			return "", ""
		}
		return key, value
	}
}

func main() {
	storage = make(map[string]string)
	listener, err := net.Listen("tcp", addr)
    if err != nil {
		panic(err) 
	} else {
		fmt.Printf("Runnig key-value server on port %v \n", addr)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Server error: failed to accept connection %v\n", err)
			continue
		}
		fmt.Printf("accepted connection from %v\n", conn.RemoteAddr())
		handle(conn) 
	}
	listener.Close()
}