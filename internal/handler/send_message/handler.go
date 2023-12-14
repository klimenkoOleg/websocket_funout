package items_delete_inactive

import (
	// "context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/klimenkoOleg/websocket_funout/internal/message"
	"github.com/klimenkoOleg/websocket_funout/internal/storage"
)

type Handler struct {
	storage  *storage.DeviceStorage
	dispatch chan message.Message
}

func New(storage *storage.DeviceStorage, dispatcher chan message.Message) *Handler {
	return &Handler{
		storage:  storage,
		dispatch: dispatcher,
	}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var msg message.Message
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid JSON message format: %v", err), http.StatusBadRequest)
		return
	}

	h.dispatch <- msg

	writer.WriteHeader(http.StatusOK)
}
