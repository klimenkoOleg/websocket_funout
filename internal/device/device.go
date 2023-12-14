package device

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/klimenkoOleg/websocket_funout/internal/message"
)

type Device struct {
	ID     uuid.UUID
	Conn   *websocket.Conn
	Logger Logger
}

func (d *Device) Send(msg message.Message) error {
	// clientsMux.Lock()
	// defer clientsMux.Unlock()

	// if client, ok := clients[deviceID]; ok {
	err := d.Conn.WriteJSON(msg)
	if err != nil {
		return fmt.Errorf("Connection error: %w", err)
	}
	// } else {
	// 	log.Printf("Client with device ID %s not found\n", deviceID)
	// }

	return nil
}
