package handle_devices

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/klimenkoOleg/websocket_funout/internal/device"
	"github.com/klimenkoOleg/websocket_funout/internal/storage"
)

type Handler struct {
	storage  *storage.DeviceStorage
	upgrader *websocket.Upgrader
	logger   Logger
}

func New(storage *storage.DeviceStorage, logger Logger) *Handler {
	return &Handler{
		logger:  logger,
		storage: storage,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	conn, err := h.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		h.logger.Error(err)
		return
	}

	defer conn.Close() // ignore possible closing errors

	// var clientID string
	params := request.URL.Query()
	deviceID := params.Get("device_id")

	id, err := uuid.Parse(deviceID)
	if err != nil {
		h.logger.Error("device_id is not UUID", err)
	}

	device := &device.Device{ID: id, Conn: conn, Logger: h.logger}

	// if deviceID != "" {
	// 	clientID = deviceID
	// } else {
	// 	clientID = "broadcast"
	// }

	h.storage.Store(device)

	defer func() {
		h.storage.Delete(device.ID) // todo try to disconnect
	}()

	for {
		// TODO exit
		_, _, err := conn.ReadMessage()
		// var msg message.Message
		// err := conn.ReadJSON(&msg)
		if err != nil {
			h.logger.Warn("client disconnected, id=%v", err)
			break
		}

		// if msg.DeviceID == "" {
		// 	broadcastMessage(msg)
		// } else {
		// 	sendMessageToClient(msg.DeviceID, msg)
		// }
		time.Sleep(time.Second)
		h.logger.Debug("listening, number of clients: ", h.storage.Count())
	}

}
