package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {

	//create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("redis clone server running and waiting for connections")

	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid command:", command)

		}

		_ = handler(args)

	})

	if err != nil {
		return
	}

	//Listen for connection
	conn, err := l.Accept()
	if err != nil {
		fmt.Println()
		return
	}

	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(conn)

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid command:", command)
			err := writer.Write(Value{typ: "string", str: ""})
			if err != nil {
				return
			}
			continue
		}

		if command == "SET" {
			err := aof.Write(value)
			if err != nil {
				return
			}
		}

		result := handler(args)
		err = writer.Write(result)
		if err != nil {
			return
		}
	}

}
