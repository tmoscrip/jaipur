package server

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/tmoscrip/jaipur/internal/game"
)

type Server struct {
	activeGame game.Game
	address    string
}

func (server Server) startTCPServer(done chan bool) {
	listener, err := net.Listen("tcp", server.address)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()
	done <- true

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}
		go server.handleConnection(conn)
	}
}

func (server Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	for {
		var command GameCommandWrapper
		if err := decoder.Decode(&command); err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Client closed the connection.")
				return
			}
			fmt.Println("Error decoding gob:", err)
			return
		}

		fmt.Printf("Received command: %s from player %d for game %d\n", command.Command.Action, command.Player, command.GameId)
		newGame, err := command.Run(server.activeGame)
		if err != nil {
			fmt.Println("Error running command:", err)
			return
		}

		server.activeGame = newGame

		if err := encoder.Encode(newGame); err != nil {
			fmt.Println("Error encoding response:", err)
			return
		}
	}
}
