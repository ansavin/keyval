package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var storage map[string]string

var addr = ":1234"
var debug = true

type writerToConnection interface {
	WriteString(string) (int, error)
	Flush() error
}

type connectionReader interface {
}

func listenTo(conn net.Conn) error {

	defer func() {
		fmt.Printf("Closing connection from %v \n", conn.RemoteAddr())
		conn.Close()
	}()
	writer := bufio.NewWriter(conn)
loop:
	for {
		data, err := waitForRequest(conn)
		if err != nil {
			sendResponce(writer, err.Error())
		}
		if len(data) < 1 {
			sendResponce(writer, "No command was specified")
			break loop
		}
		if debug {
			fmt.Println("Handling request", data)
		}
		switch command := data[0]; command {
		case "set":
			uploadToStorage(data[1:])
		case "show":
			printStorageData(writer)
		case "help":
			printHelp(writer)
		case "quit":
			break loop
		case "get":
			getKey(writer, data[1:])
		case "del":
			deleteKey(data)
		default:
			sendResponce(writer, fmt.Sprintf("Wrong command: %v \n", command))
			break loop
		}
	}
	return nil
}

func waitForRequest(conn net.Conn) ([]string, error) {
	reader := bufio.NewReader(conn)
	scaner := bufio.NewScanner(reader)
	scanned := scaner.Scan()
	if !scanned {
		if err := scaner.Err(); err != nil {
			return nil, fmt.Errorf("%v(%v)", err, conn.RemoteAddr())
		}
	}
	data := strings.Fields(scaner.Text())
	return data, nil
}

func sendResponce(writer writerToConnection, message string) {
	writer.WriteString(message)
	writer.Flush()
}

func printHelp(writer writerToConnection) {

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
	sendResponce(writer, message)
}
func printStorageData(writer writerToConnection) {
	var message string
	for key, value := range storage {
		message += fmt.Sprintf("%v: %v\n", key, value)
	}
	sendResponce(writer, message)
}

func deleteKey(data []string) error {
	key, value, err := findDataByKey(data[1:])
	if err != nil {
		return err
	}
	if key != "" && value != "" {
		delete(storage, key)
		return nil
	}
	return fmt.Errorf("storage key not found or not specified")
}

func uploadToStorage(data []string) {
	for index := 0; index < len(data); index++ {
		key := index
		if index++; index >= len(data) {
			storage[data[key]] = ""
		} else {
			storage[data[key]] = data[index]
		}
	}
}

func findDataByKey(data []string) (string, string, error) {
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

func getKey(writer writerToConnection, data []string) {
	key, value, err := findDataByKey(data)
	if err == nil {
		sendResponce(writer, fmt.Sprintf("%v:%v\n", key, value))
	} else {
		sendResponce(writer, fmt.Sprintf("storage key not found or not specified\n"))
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
		listenTo(connection)
	}
}
