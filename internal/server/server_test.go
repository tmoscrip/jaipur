package server

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"testing"

	"github.com/tmoscrip/jaipur/internal/game"
)

func TestServer(t *testing.T) {
	// generate random port
	port := rand.Int()%1000 + 5000
	server := Server{
		activeGame: game.NewGame(),
		address:    fmt.Sprintf("localhost:%d", port),
	}
	done := make(chan bool)
	go server.startTCPServer(done)
	<-done

	conn, err := net.Dial("tcp", server.address)
	if err != nil {
		t.Errorf("Error dialing server: %v", err)
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	command := NewSetPlayerNameCommand(0, "test")
	err = encoder.Encode(command)
	if err != nil {
		t.Errorf("Error encoding command: %v", err)
	}

	var newGame game.Game
	err = decoder.Decode(&newGame)
	if err != nil {
		t.Errorf("Error decoding game: %v", err)
	}

	if newGame.Players.Get(0).Name != "test" {
		t.Errorf("Expected name to be 'test', got %v", newGame.Players.Get(0).Name)
	}
}
