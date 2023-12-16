//go:generate mockgen -source $GOFILE -destination mock_test.go -package ${GOPACKAGE}
package send_message

import (
	"github.com/google/uuid"

	"github.com/klimenkoOleg/websocket_funout/internal/device"
)

type Logger interface {
	Debug(args ...interface{})
	Error(args ...interface{})
}

type Storage interface {
	Count() int
	IsDeviceRegistered(id uuid.UUID) bool
	Store(d *device.Device)
	Delete(id uuid.UUID)
}
