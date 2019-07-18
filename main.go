package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var storage map[string]string

var addr string = ":1234"
var debug bool = true

type writer_to_connection interface {
	WriteString(string) (int, error)
	Flush() error
}

type connection_reader interface {
}

func listen_to(conn net.Conn) error {

	defer func() {
		fmt.Printf("Closing connection from %v \n", conn.RemoteAddr())
		conn.Close()
	}()
	writer := bufio.NewWriter(conn)
loop:
	for {
		data, err := wait_foe_request(conn)
		if err != nil {
			send_responce(writer, err.Error())
		}
		if len(data) < 1 {
			send_responce(writer, "No command was specified")
			break loop
		}
		if debug {
			fmt.Println("Handling request", data)
		}
		switch command := data[0]; command {
		case "set":
			upload_to_storage(data[1:])
		case "show":
			print_storage_data(writer)
		case "help":
			print_help(writer)
		case "quit":
			break loop
		case "get":
			get_key(writer, data[1:])
		case "del":
			delete_key(data)
		default:
			send_responce(writer, fmt.Sprintf("Wrong command: %v \n", command))
			break loop
		}
	}
	return nil
}

func wait_foe_request(conn net.Conn) ([]string, error) {
	reader := bufio.NewReader(conn)
	scaner := bufio.NewScanner(reader)
	scanned := scaner.Scan()
	if !scanned {
		if err := scaner.Err(); err != nil {
			return nil, fmt.Errorf("%v(%v)\n", err, conn.RemoteAddr())
		}
	}
	data := strings.Fields(scaner.Text())
	return data, nil
}

func send_responce(writer writer_to_connection, message string) {
	writer.WriteString(message)
	writer.Flush()
}

func print_help(writer writer_to_connection) {

	message := `
	Usage: 
	set <key> <value>
	get <key>
	del <key>
	show
	help
	quit
	`
	message += "\n"
	send_responce(writer, message)
}
func print_storage_data(writer writer_to_connection) {
	var message string
	for key, value := range storage {
		message += fmt.Sprintf("%v: %v\n", key, value)
	}
	send_responce(writer, message)
}

func delete_key(data []string) error {
	key, value, err := find_data_by_key(data[1:])
	if err != nil {
		return err
	}
	if key != "" && value != "" {
		delete(storage, key)
		return nil
	} else {
		return fmt.Errorf("storage key not found or not specified")
	}
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

func find_data_by_key(data []string) (string, string, error) {
	if len(data) < 1 {
		return "", "", fmt.Errorf("Server error: no storage key specified")
	} else {
		key := data[0]
		value, ok := storage[key]
		if !ok {
			return "", "", fmt.Errorf("Server error: wrong storage key: %v", key)
		}
		return key, value, nil
	}
}

func get_key(writer writer_to_connection, data []string) {
	key, value, err := find_data_by_key(data)
	if err == nil {
		send_responce(writer, fmt.Sprintf("%v: %v\n", key, value))
	} else {
		send_responce(writer, fmt.Sprintf("storage key not found or not specified\n"))
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
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Server error: failed to accept connection %v\n", err)
			continue
		}
		fmt.Printf("accepted connection from %v\n", connection.RemoteAddr())
		listen_to(connection)
	}
	listener.Close()
}
