package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/klimenkoOleg/websocket_funout/internal/dto"
	"github.com/klimenkoOleg/websocket_funout/internal/handler/handle_devices"
	"github.com/klimenkoOleg/websocket_funout/internal/handler/send_message"
	"github.com/klimenkoOleg/websocket_funout/internal/infra/logger"
	"github.com/klimenkoOleg/websocket_funout/internal/server"
	"github.com/klimenkoOleg/websocket_funout/internal/storage"
)

func main() {
	var err error
	logger := logger.MustInitLogger()
	defer logger.Sync() // flush buffer

	logger.Debug()

	defer func() {
		if panicErr := recover(); panicErr != nil {
			logger.Error("recover", zap.Reflect("recover error", panicErr))
			os.Exit(1)
		}

		if err != nil {
			logger.Error("error left", zap.Error(err))
			os.Exit(1)
		}
	}()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	dispatch := make(chan dto.Message)

	deviceStorage := storage.New(dispatch, logger)

	go deviceStorage.Start(ctx)

	sendMessageHandler := send_message.New(deviceStorage, dispatch, logger)
	devicesHandler := handle_devices.New(deviceStorage, logger, &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})

	mux := http.DefaultServeMux
	mux.HandleFunc("/send", sendMessageHandler.ServeHTTP)
	mux.Handle("/ws", devicesHandler)

	serverListener := server.New(
		server.WithLogger(logger),
		server.WithShutdownTimeout(500*time.Millisecond),
		server.WithOnShutdown(func() {
			deviceStorage.Stop()
		}),
	)

	logger.Fatal(serverListener.Listen(ctx, mux))
}
