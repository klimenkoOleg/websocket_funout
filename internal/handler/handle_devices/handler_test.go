package handle_devices

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Initialize the components with mock implementations
	mockStorage := NewMockStorage(ctrl)
	mockLogger := NewMockLogger(ctrl)
	// mockConnector := NewMockConn(ctrl)
	// mockUpgrader := NewMockUpgrader(ctrl)

	// mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConnector, nil)
	mockStorage.EXPECT().IsDeviceRegistered(gomock.Any()).Return(false).Times(1)
	mockStorage.EXPECT().Store(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Debug(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()
	mockStorage.EXPECT().Delete(gomock.Any()).AnyTimes()

	responseRecorder := httptest.NewRecorder()

	handler := New(mockStorage, mockLogger, &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})

	s := httptest.NewServer(http.HandlerFunc(handler.ServeHTTP))
	defer s.Close()
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/?device_id=" + uuid.New().String()
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	// Assert the success conditions
	assert.Equal(t, http.StatusOK, responseRecorder.Code, "Unexpected HTTP status code")
}
