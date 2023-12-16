package handle_devices

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/klimenkoOleg/websocket_funout/internal/device"
	"github.com/klimenkoOleg/websocket_funout/internal/dto"
)

type Handler struct {
	storage  Storage
	upgrader Upgrader
	logger   Logger
}

func New(storage Storage, logger Logger, upgrader Upgrader) *Handler {
	return &Handler{
		logger:   logger,
		storage:  storage,
		upgrader: upgrader,
	}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	conn, err := h.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		h.logger.Error(err)
		return
	}

	params := request.URL.Query()
	deviceID := params.Get("device_id")

	id, err := uuid.Parse(deviceID)
	if err != nil {
		h.logger.Error("device_id is not UUID", err)
		return
	}

	if h.storage.IsDeviceRegistered(id) {
		h.logger.Debug("attempt to register duplicated device_id=", id)
		conn.WriteJSON(dto.Error{400, "the device with such device_id is already registered"})
		conn.Close()
		return
	}

	device := device.New(id, conn)

	h.storage.Store(device)
	defer func(id uuid.UUID) {
		h.logger.Debug("disconnecting with device id=", id)
		h.storage.Delete(device.Id)
		conn.Close()
	}(id)

	h.logger.Debug("opened connection with device id=", device.Id)

	// no need in `select` or additional quit channel here: conn.Close() made elsewhere also terminates conn.ReadMessage()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			h.logger.Warn("client disconnected, id=", err)
			break
		}
		h.logger.Debug("listening, number of clients: ", h.storage.Count())
	}
}
