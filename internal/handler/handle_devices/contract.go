//go:generate mockgen -source $GOFILE -destination mock_test.go -package ${GOPACKAGE}
package handle_devices

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/klimenkoOleg/websocket_funout/internal/device"
)

type Logger interface {
	Debug(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
}

type Storage interface {
	Count() int
	IsDeviceRegistered(id uuid.UUID) bool
	Store(d *device.Device)
	Delete(id uuid.UUID)
}

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)
}
