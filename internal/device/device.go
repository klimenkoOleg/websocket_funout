package device

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/klimenkoOleg/websocket_funout/internal/dto"
)

func New(
	id uuid.UUID,
	conn Connector,
) *Device {
	return &Device{Id: id, conn: conn}
}

type Device struct {
	Id   uuid.UUID
	conn Connector
}

// Disconnect  Closes connection, flushed buffers
func (d *Device) Disconnect() error {
	return d.conn.Close()
}

// Send closes connection on failure to send.
func (d *Device) Send(msg dto.Message) error {
	err := d.conn.WriteJSON(msg)
	if err != nil {
		if disconnectErr := d.Disconnect(); disconnectErr != nil {
			err = fmt.Errorf("error closing device connection: %w", disconnectErr)
		}
		return fmt.Errorf("Connection error: %w", err)
	}
	return nil
}
