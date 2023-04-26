package phocus_sensors

import (
	"fmt"
	"log"
	"time"

	"github.com/wolffshots/ha_types/device_classes"
	"github.com/wolffshots/ha_types/state_classes"
	"github.com/wolffshots/ha_types/units"
	"github.com/wolffshots/phocus_mqtt"
)

// Sensor is the shape of the sensor for the MQTT Home Assistant integration
type Sensor struct {
	SensorTopic   string                     // "homeassistant/sensor/phocus/start_time/config" must end in /config
	UniqueId      string                     // "unique_id": "phocus_qpgs1_ac_output_apparent_power",
	Unit          units.Unit                 // "unit_of_measurement": "VA",
	StateClass    state_classes.StateClass   // "state_class": "measurement",
	DeviceClass   device_classes.DeviceClass // "device_class": "apparent_power",
	Name          string                     // "name": "QPGS1 AC Output Apparent Power",
	ValueTemplate string                     // "value_template": "{{ value_json.ACOutputApparentPower }}",
	StateTopic    string                     // "state_topic": "phocus/stats/qpgs1",
	Icon          string                     // "icon": "mdi:battery",
}

var sensors = []Sensor{
	{
		SensorTopic:   "homeassistant/sensor/phocus/start_time/config",
		UniqueId:      "phocus_start_time",
		Unit:          units.None,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Timestamp,
		Name:          "Start Time",
		ValueTemplate: "",
		StateTopic:    "phocus/stats/start_time",
		Icon:          "mdi:clock",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/error/config",
		UniqueId:      "phocus_last_error",
		Unit:          units.None,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.None,
		Name:          "Last Reported Error",
		ValueTemplate: "",
		StateTopic:    "phocus/stats/error",
		Icon:          "mdi:hammer-wrench",
	},
}

// Register adds some sensors to Home Assistant MQTT
func Register() error {
	log.Println("Registering sensors")
	for _, input := range sensors {
		log.Printf("Registering %s\n", input.Name)

		sensor_string := fmt.Sprintf(
			"{\""+
				"unique_id\":\"%s\",\""+
				"name\":\"%s\",\""+
				"state_topic\":\"%s\",\""+
				"icon\":\"%s\",\""+
				"unit\":\"%s\",\""+
				"state_class\":\"%s\",\""+
				"device_class\":\"%s\",\""+
				"device\":{\"name\":\"phocus\",\""+
				"identifiers\":[\"phocus\"],\""+
				"model\":\"phocus\",\""+
				"manufacturer\":\"phocus\",\""+
				"sw_version\":\"1.1.0\"},\""+
				"force_update\":false",
			input.UniqueId,
			input.Name,
			input.StateTopic,
			input.Icon,
			input.Unit,
			input.StateClass,
			input.DeviceClass,
		)

		if input.ValueTemplate != "" {
			sensor_string += "value_template\":\"%s\",\""
		}
		sensor_string += "}"

		err := phocus_mqtt.Send(input.SensorTopic, 0, true, sensor_string, 10)
		if err != nil {
			log.Printf("Failed to send initial setup stats to MQTT with err: %v", err)
			return err
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}
