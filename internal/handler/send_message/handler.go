package send_message

import (
	// "context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/klimenkoOleg/websocket_funout/internal/dto"
)

type Handler struct {
	storage  Storage
	dispatch chan dto.Message
	logger   Logger
}

func New(storage Storage, dispatcher chan dto.Message, logger Logger) *Handler {
	return &Handler{
		storage:  storage,
		dispatch: dispatcher,
		logger:   logger,
	}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	var msg dto.Message
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid JSON dto format: %v", err), http.StatusBadRequest)

		return
	}

	if msg.DeviceID != nil && !h.storage.IsDeviceRegistered(*msg.DeviceID) {
		http.Error(writer, "Requested device not connected", http.StatusNotFound)
		writer.WriteHeader(http.StatusOK)

		return
	}

	h.logger.Debug(fmt.Sprintf("Processing incoming dto: %+v", msg))

	h.dispatch <- msg

	writer.WriteHeader(http.StatusOK)
}
