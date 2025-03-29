package main

import (
	"fmt"
	"net"
)

func main() {

	//create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
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

		_ = value

		writer := NewWriter(conn)
		err = writer.Write(Value{typ: "string", str: "OK"})
		if err != nil {
			return
		}

	}

}
