package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/klimenkoOleg/websocket_funout/internal/device"
	"github.com/klimenkoOleg/websocket_funout/internal/dto"
)

type DeviceStorage struct {
	devices map[uuid.UUID]*device.Device
	// register   chan *device.Device
	// deregister chan *device.Device
	// quit       chan struct{}
	devicesMux sync.Mutex
	dispatch   chan dto.Message
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

func New(dispatch chan dto.Message, logger Logger) *DeviceStorage {
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

func (ds *DeviceStorage) IsDeviceRegistered(id uuid.UUID) bool {
	_, ok := ds.devices[id]
	return ok
}

func (ds *DeviceStorage) Stop() {
	ds.logger.Debug("stopping storage...")
	ds.closeDevicesConnections()
}

func (ds *DeviceStorage) Start(ctx context.Context) {
	func() {
		for {
			select {
			// case <-ds.quit:
			// 	ds.logger.Println("quitting...")
			// 	return
			case <-ctx.Done():
				ds.logger.Debug("shutting down storage... # of devices+", len(ds.devices))
				for _, d := range ds.devices {
					d.Disconnect()
				}
				// ds.closeDevicesConnections()
				return
			case msg := <-ds.dispatch:
				if msg.DeviceID == nil {
					ds.logger.Debug(fmt.Sprintf("broadcast dto: %+v", msg))
					ds.broadcast(msg)
				} else {
					ds.logger.Debug(fmt.Sprintf("sending to one device, dto: %+v", msg))
					ds.sendToDevice(*msg.DeviceID, msg)
				}
				// for _, d := range ds.devices {
				// 	d.send(dto)
				// }
				// case device := <-ds.register:
				// 	ds.Store(device)
				// case device := <-ds.deregister:
				// 	ds.Delete(device.id)
			}
		}
	}()
}

func (ds *DeviceStorage) closeDevicesConnections() {
	ds.logger.Debug("closeDevicesConnections, # of devices: ", len(ds.devices))
	for _, d := range ds.devices {
		d.Disconnect()
	}
}

// todo send the func, not hardcode
func (ds *DeviceStorage) broadcast(msg dto.Message) {
	// This is a place for performance improvement: run Send in a goroutine for each device.
	// But I prefer not to do premature optimization.
	for _, d := range ds.devices {
		err := d.Send(msg)
		if err != nil {
			ds.logger.Error("client sending error", err)
			ds.Delete(d.Id)
		}
	}
}

func (ds *DeviceStorage) sendToDevice(deviceId uuid.UUID, msg dto.Message) {
	device, ok := ds.devices[deviceId]
	if !ok {
		ds.logger.Warn("device not found by deviceId: ", deviceId)
		return
	}
	// this is a place for performance improvement: run Send in a goroutine for each device
	err := device.Send(msg)
	if err != nil {
		delete(ds.devices, device.Id)
	}
}

// Performance hint: mutex could be replaced by a Goroutine but I'd love to aboit premature optimization.
func (ds *DeviceStorage) Store(d *device.Device) {
	ds.logger.Debug("adding device #", d.Id)

	ds.devicesMux.Lock()
	ds.devices[d.Id] = d
	ds.devicesMux.Unlock()
}

func (ds *DeviceStorage) Delete(id uuid.UUID) {
	ds.logger.Debug("deleting device #", id)
	ds.devicesMux.Lock()
	delete(ds.devices, id)
	ds.devicesMux.Unlock()
}
