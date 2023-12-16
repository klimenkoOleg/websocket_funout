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
	devices  map[uuid.UUID]*device.Device
	store    chan *device.Device
	delete   chan uuid.UUID
	dispatch chan dto.Message
	logger   Logger
	mu       sync.Mutex
}

func New(dispatch chan dto.Message, logger Logger) *DeviceStorage {
	deviceStorage := &DeviceStorage{
		devices:  make(map[uuid.UUID]*device.Device),
		store:    make(chan *device.Device),
		delete:   make(chan uuid.UUID),
		dispatch: dispatch,
		logger:   logger,
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
			case <-ctx.Done():
				ds.logger.Debug("shutting down storage... # of devices+", len(ds.devices))

				for _, d := range ds.devices {
					ds.logger.Debug("closed devices connection #", d.Id)

					if err := d.Disconnect(); err != nil {
						ds.logger.Warn("error closing device connection: ", err)
					}
				}
				return
			case msg := <-ds.dispatch:
				if msg.DeviceID == nil {
					ds.logger.Debug(fmt.Sprintf("broadcast dto: %+v", msg))
					ds.broadcast(msg)
				} else {
					ds.logger.Debug(fmt.Sprintf("sending to one device, dto: %+v", msg))
					ds.sendToDevice(*msg.DeviceID, msg)
				}
			case device := <-ds.store:
				// no need to sync.Mutex.Lock()/Unlock - this is single thread changing the map, and select processes only one case at a time
				ds.devices[device.Id] = device
			case deviceId := <-ds.delete:
				// no need to sync.Mutex.Lock()/Unlock - this is single thread changing the map
				delete(ds.devices, deviceId)
			}
		}
	}()
}

func (ds *DeviceStorage) closeDevicesConnections() {
	ds.logger.Debug("closeDevicesConnections, # of devices: ", len(ds.devices))

	for _, d := range ds.devices {
		ds.logger.Debug("closed devices connection #", d.Id)

		if err := d.Disconnect(); err != nil {
			ds.logger.Warn("error closing device connection: ", err)
		}
	}
}

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
	err := device.Send(msg)
	if err != nil {
		delete(ds.devices, device.Id)
	}
}

func (ds *DeviceStorage) Store(d *device.Device) {
	ds.logger.Debug("adding device #", d.Id)
	ds.store <- d
}

func (ds *DeviceStorage) Delete(id uuid.UUID) {
	ds.logger.Debug("deleting device #", id)
	ds.delete <- id
}
