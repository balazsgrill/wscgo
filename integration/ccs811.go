package integration

import (
	"github.com/home2mqtt/wscgo/config"
	"github.com/home2mqtt/wscgo/devices"
	"github.com/home2mqtt/wscgo/protocol"
)

type ccs811ConfigPartParser struct{}

type ccs811SensorConfiguration struct {
	protocol.BasicDeviceConfig
	Name      string `ini:"name,omitempty"`
	ECO2Topic string `ini:"eco2topic,omitempty"`
	VOCTopic  string `ini:"voctopic,omitempty"`
}

func (*ccs811ConfigPartParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	id := section.GetID()
	c := &ccs811SensorConfiguration{
		BasicDeviceConfig: protocol.BasicDeviceConfig{
			Name:     id,
			ObjectId: id,
		},
	}
	section.FillData(c)
	d := &devices.CCS811Config{
		Address:  0x5A,
		Duration: 10,
	}
	section.FillData(d)

	eco2conf := &protocol.SensorConfig{
		BasicDeviceConfig: protocol.BasicDeviceConfig{
			Name:     c.Name + " eCO2",
			ObjectId: id + "_eco2",
		},
		Icon:              "mdi:periodic-table-co2",
		Topic:             c.ECO2Topic,
		UnitOfMeasurement: "ppm",
	}
	vocconf := &protocol.SensorConfig{
		BasicDeviceConfig: protocol.BasicDeviceConfig{
			Name:     c.Name + " VOC",
			ObjectId: id + "_voc",
		},
		Icon:              "mdi:weather-windy",
		Topic:             c.VOCTopic,
		UnitOfMeasurement: "ppb",
	}
	context.AddDeviceInitializer(config.SLDevice, func(context config.RuntimeContext) error {
		device, err := devices.CreateCCS811(d)
		if err != nil {
			return err
		}
		context.AddDevice(device)
		context.AddProtocol(protocol.IntegrateSensor(device.ECO2(), eco2conf))
		context.AddProtocol(protocol.IntegrateSensor(device.VOC(), vocconf))
		return nil
	})
	return nil
}

func init() {
	config.RegisterConfigurationPartParser("ccs811", &ccs811ConfigPartParser{})
}
