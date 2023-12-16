package send_message

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/klimenkoOleg/websocket_funout/internal/dto"
)

func TestHandler_ServeHTTP_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)

	// Initialize the components with mock implementations
	mockLogger := NewMockLogger(ctrl)

	mockLogger.EXPECT().Debug(gomock.Any()).AnyTimes()

	dispatchChannel := make(chan dto.Message, 1)
	handler := New(mockStorage, dispatchChannel, mockLogger)

	// Create a test request
	message := dto.Message{
		ID:      uuid.New(),
		Message: "Test message content",
		Kind:    123456789,
	}
	jsonData, _ := json.Marshal(message)
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
	responseRecorder := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(responseRecorder, request)

	// Assert the success conditions
	assert.Empty(t, responseRecorder.Body.String(), "Unexpected response body")
	assert.Equal(t, http.StatusOK, responseRecorder.Code, "Unexpected HTTP status code")
	assert.Equal(t, 1, len(dispatchChannel), "Message should be sent to the dispatch channel")
	assert.Equal(t, message, <-dispatchChannel, "Sent message should match the expected message")
}
