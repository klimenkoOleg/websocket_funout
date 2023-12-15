package device

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/klimenkoOleg/websocket_funout/internal/dto"
)

func New(
	id uuid.UUID,
	conn *websocket.Conn,
	logger Logger,
) *Device {
	return &Device{Id: id, conn: conn, logger: logger}
}

type Device struct {
	Id     uuid.UUID
	conn   *websocket.Conn
	logger Logger
}

// Disconnect  Closes connection, flushed buffers
func (d *Device) Disconnect() {
	err := d.conn.Close()
	d.logger.Debug("closed devices connection #", d.Id)
	// we're exiting, so aren't caring about sending errors up, just logging it
	if err != nil {
		d.logger.Warn("error closing device connection: ", err)
	}
}

// Send closes connection on failure to send.
func (d *Device) Send(msg dto.Message) error {
	err := d.conn.WriteJSON(msg)
	if err != nil {
		d.Disconnect()
		return fmt.Errorf("Connection error: %w", err)
	}
	return nil
}
