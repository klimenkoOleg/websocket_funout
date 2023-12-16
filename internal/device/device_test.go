package device

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/klimenkoOleg/websocket_funout/internal/dto"
)

func TestDevice_Send_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deviceID := uuid.New()

	mockConnector := NewMockConnector(ctrl)
	mockConnector.EXPECT().WriteJSON(gomock.Any()).Times(1)

	device := New(deviceID, mockConnector)

	message := dto.Message{ID: uuid.New(), Kind: 1, Message: "Test Message"}
	err := device.Send(message)

	require.NoError(t, err)
}

func TestDevice_Disconnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deviceID := uuid.New()
	mockConnector := NewMockConnector(ctrl)
	mockConnector.EXPECT().Close().Times(1)

	device := New(deviceID, mockConnector)

	err := device.Disconnect()

	require.NoError(t, err)
}
