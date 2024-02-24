package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/LasithaPrabodha/redis-like-server/internal"
)

func Start(address string) {

	// Create a new server
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := internal.NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	fmt.Println("Server started. Listening on", address)

	aof.Read(func(value internal.Value) {
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := internal.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, aof)
	}
}

func handleConnection(conn net.Conn, aof *internal.Aof) {
	defer conn.Close()

	resp := internal.NewResp(conn)
	for {
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.Typ != "array" || len(value.Array) == 0 {
			fmt.Println("Invalid request")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := internal.NewWriter(conn)

		handler, ok := internal.Handlers[command]
		if !ok {
			fmt.Println("Invalid command:", command)
			writer.Write(internal.Value{Typ: "string", Str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			if err := aof.Write(value); err != nil {
				fmt.Println("Failed to write to AOF:", err)
				continue
			}
		}

		result := handler(args)
		writer.Write(result)
	}
}
