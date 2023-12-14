package storage

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/klimenkoOleg/websocket_funout/internal/device"
	"github.com/klimenkoOleg/websocket_funout/internal/message"
)

type DeviceStorage struct {
	devices map[uuid.UUID]*device.Device
	// register   chan *device.Device
	// deregister chan *device.Device
	// quit       chan struct{}
	devicesMux sync.Mutex
	dispatch   chan message.Message
	logger     Logger

	// func NewHub(ch chan counter, quit chan struct{}) *hub {
	// return &hub{
	// clients:    make(map[int]*client),
	// register:   make(chan *client),
	// deregister: make(chan *client),
	// quit:       quit,
	// dispatch:   ch,
}

// }
// }

func New(dispatch chan message.Message, logger Logger) *DeviceStorage {
	deviceStorage := &DeviceStorage{
		devices: make(map[uuid.UUID]*device.Device),
		// register:   make(chan *device.Device),
		// deregister: make(chan *device.Device),
		dispatch: dispatch,
		// quit:       quit,
		logger: logger,
	}

	return deviceStorage
}

func (ds *DeviceStorage) Count() int {
	return len(ds.devices)
}

func (ds *DeviceStorage) Start(ctx context.Context) {
	go func() {
		for {
			select {
			// case <-ds.quit:
			// 	ds.logger.Println("quitting...")
			// 	return
			case <-ctx.Done():
				ds.logger.Debug("quitting...")
				return
			case msg := <-ds.dispatch:
				if msg.DeviceID == nil {
					ds.broadcast(msg)
				} else {
					ds.sendToDevice(*msg.DeviceID, msg)
				}
				// for _, d := range ds.devices {
				// 	d.send(message)
				// }
				// case device := <-ds.register:
				// 	ds.Store(device)
				// case device := <-ds.deregister:
				// 	ds.Delete(device.ID)
			}
		}
	}()
}

func (ds *DeviceStorage) broadcast(msg message.Message) {
	for _, d := range ds.devices {
		err := d.Send(msg)
		if err != nil {
			ds.logger.Error("client sending error", err)
			ds.Delete(d.ID)
		}
	}
}

func (ds *DeviceStorage) sendToDevice(deviceId uuid.UUID, msg message.Message) {
	device, ok := ds.devices[deviceId]
	if !ok {
		ds.logger.Warn("device not found by deviceId: ", deviceId)
		return
	}
	device.Send(msg)
}

func (ds *DeviceStorage) Store(d *device.Device) {
	ds.devicesMux.Lock()
	ds.devices[d.ID] = d
	ds.devicesMux.Unlock()
}

func (ds *DeviceStorage) Delete(id uuid.UUID) {
	ds.devicesMux.Lock()
	delete(ds.devices, id)
	ds.devicesMux.Unlock()
}
