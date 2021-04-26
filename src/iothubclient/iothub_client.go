package iothubclient

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/amenzhinsky/iothub/iotdevice"
	iotmqtt "github.com/amenzhinsky/iothub/iotdevice/transport/mqtt"
)

// DeviceInfo is the type for passing the device info tuples
// It is a hierarcical map with key/value pairs
type DeviceInfo map[string]interface{}

// New creates a new Iotdevice client from the device connection string cs
func New(cs string) (*iotdevice.Client, error) {
	c, err := iotdevice.NewFromConnectionString(
		// <transport>, <connection string>,
		iotmqtt.New(), cs)
	return c, err
}

// SetStaticDeviceInfo writes device info this iotedge module's twin
// d is a hierarcical map which is placed currently unter the "reported/versions"
// properties within the module twin
func SetStaticDeviceInfo(c *iotdevice.Client, d DeviceInfo) error {
	if d == nil {
		return errors.New("DeviceInfo is nil")
	}

	// connect to the iothub
	if err := waitForIotHubConnection(c); err != nil {
		return err
	}

	s := makeStaticDeviceInfo(d)
	if _, err := c.UpdateTwinState(context.Background(), s); err != nil {
		return err
	}
	return nil
}

// Wait forever for IotHub Connection
func waitForIotHubConnection(c *iotdevice.Client) error {
	for {
		err := c.Connect(context.Background())

		if err == nil {
			fmt.Println("connect to iothub ok")
			return nil
		}
		time.Sleep(time.Second)
	}
}

func makeStaticDeviceInfo(d DeviceInfo) iotdevice.TwinState {
	s := make(iotdevice.TwinState)

	s["versions"] = d
	return s
}
