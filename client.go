package ets2

import "log"
import "net/http"
import "io/ioutil"
import "encoding/json"

type Client struct {
	BaseUrl string
}

type PlacementData struct {
	X, Y, Z float32
	Heading float32
	Pitch   float32
	Roll    float32
}

type VectorData struct {
	X, Y, Z float32
}

type Telemetry struct {
	Game struct {
		Connected              bool
		Paused                 bool
		GameName               string
		Time                   string
		TimeScale              float32
		NextRestStopTime       string
		Version                string
		TelemetryPluginVersion string
	}

	Truck struct {
		Id    string // Brand Id of the current truck. Standard values are: "daf", "iveco", "man", "mercedes", "renault", "scania", "volvo"
		Make  string // Localized brand name of the current truck for display purposes
		Model string // Localized model name of the current truck

		Speed              float32 // Speed in kmh
		CruiseControlSpeed float32 // Speed selected for the cruise control in km/h
		CruiseControlOn    bool
		Odometer           float32 // The value of the truck's odometer in km

		Gear          int    // Gear that is currently selected in the engine (physical gear). Positive values reflect forward gears, negative - reverse
		DisplayedGear int    // Gear that is currently displayed on the main dashboard inside the game. Positive values reflect forward gears, negative - reverse
		ForwardGears  int    // Number of forward gears on undamaged truck
		ReverseGears  int    // Number of reverse gears on undamaged truck
		ShifterType   string // Type of the shifter selected in the game's settings. One of the following values: "arcade", "automatic", "manual", "hshifter"

		EngineRpm    float32 // Current RPM value of the truck's engine (rotates per minute).
		EngineRpmMax float32 // Maximal RPM value of the truck's engine

		Fuel                   float32 // Current amount of fuel in liters
		FuelCapacity           float32 // Fuel tank capacity in litres
		FuelAverageConsumption float32 // Average consumption of the fuel in liters/km
		FuelWarningFactor      float32 // Fraction of the fuel capacity bellow which is activated the fuel warning
		FuelWarningOn          bool    // Indicates whether low fuel warning is active or not

		WearEngine       float32 // Current level of truck's engine wear/damage between 0 (min) and 1 (max)
		WearTransmission float32 // Current level of truck's transmission wear/damage between 0 (min) and 1 (max)
		WearCabin        float32 // Current level of truck's cabin wear/damage between 0 (min) and 1 (max)
		WearChassis      float32 // Current level of truck's chassis wear/damage between 0 (min) and 1 (max)
		WearWheels       float32 // Current level of truck's wheel wear/damage between 0 (min) and 1 (max)

		UserSteer    float32 // Steering received from input (-1;1). Note that it is interpreted counterclockwise. If the user presses the steer right button on digital input (e.g. keyboard) this value goes immediatelly to -1.0
		UserThrottle float32 // Throttle received from input (-1;1). If the user presses the forward button on digital input (e.g. keyboard) this value goes immediatelly to 1.0
		UserBrake    float32 // Brake received from input (-1;1). If the user presses the brake button on digital input (e.g. keyboard) this value goes immediatelly to 1.0
		UserClutch   float32 // Clutch received from input (-1;1). If the user presses the clutch button on digital input (e.g. keyboard) this value goes immediatelly to 1.0

		GameSteer    float32 // Steering as used by the simulation (-1;1). Note that it is interpreted counterclockwise. Accounts for interpolation speeds and simulated counterfoces for digital inputs
		GameThrottle float32 // Throttle pedal input as used by the simulation (0;1). Accounts for the press attack curve for digital inputs or cruise-control input
		GameBrake    float32 // Brake pedal input as used by the simulation (0;1). Accounts for the press attack curve for digital inputs. Does not contain retarder, parking or motor brake
		GameClutch   float32 // Clutch pedal input as used by the simulation (0;1). Accounts for the automatic shifting or interpolation of player input

		ShifterSlot int  // Gearbox slot the h-shifter handle is currently in. 0 means that no slot is selected
		EngineOn    bool // Indicates whether the engine is currently turned on or off
		ElectricOn  bool // Indicates whether the electric is enabled or not
		WipersOn    bool // Indicates whether wipers are currently turned on or off

		RetarderBrake     int     // Current level of the retarder brake. Ranges from 0 to RetarderStepCount.
		RetarderStepCount int     // Number of steps in the retarder. Set to zero if retarder is not mounted to the truck.
		ParkBrakeOn       bool    // Is the parking brake enabled or not
		MotorBrakeOn      bool    // Is the motor brake enabled or not
		BrakeTemperature  float32 // Temperature of the brakes in degrees celsius

		Adblue                   float32 // Amount of AdBlue in liters
		AdblueCapacity           float32 // AdBlue tank capacity in litres
		AdblueAverageConsumption float32 // Average consumption of the adblue in liters/km
		AdblueWarningOn          bool    // Is the low adblue warning active or not

		AirPressure               float32 // Pressure in the brake air tank in psi
		AirPressureWarningOn      bool    // Is the air pressure warning active or not
		AirPressureWarningValue   float32 // Pressure of the air in the tank bellow which the warning activates
		AirPressureEmergencyOn    bool    // Are the emergency brakes active as result of low air pressure or not
		AirPressureEmegrencyValue float32 // Pressure of the air in the tank bellow which the emergency brakes activate

		OilTemperature          float32 // Temperature of the oil in degrees celsius
		OilPressure             float32 // Pressure of the oil in psi
		OilPressureWarningOn    bool    // Is the oil pressure warning active or not
		OilPressureWarningLevel float32 // Pressure of the oil bellow which the warning activates

		WaterTemperature             float32 // Temperature of the water in degrees celsius
		WaterTemperatureWarningOn    bool    // Is the water temperature warning active or not
		WaterTemperatureWarningLevel float32 // Temperature of the water above which the warning activates

		BatteryVoltage             float32 // Voltage of the battery in volts
		BatteryVoltageWarningOn    bool    // Is the battery voltage/not charging warning active or not
		BatteryVoltageWarningValue float32 // Voltage of the battery bellow which the warning activates

		LightsDashboardValue float32 // Intensity of the dashboard backlight between 0 (off) and 1 (max)
		LightsDashboardOn    bool    // Is the dashboard backlight currently turned on or off

		BlinkerLeftActive  bool // Indicates whether the left blinker currently emits light or not
		BlinkerRightActive bool // Indicates whether the right blinker currently emits light or not
		BlinkerLeftOn      bool // Is left blinker currently turned on or off
		BlinkerRightOn     bool // Is right blinker currently turned on or off

		LightsParkingOn  bool // Are parking lights enabled or not
		LightsBeamLowOn  bool // Are low beam lights enabled or not
		LightsBeamHighOn bool // Are high beam lights enabled or not
		LightsAuxFrontOn bool // Are auxiliary front lights active or not
		LightsAuxRoofOn  bool // Are auxiliary roof lights active or not
		LightsBeaconOn   bool // Are beacon lights enabled or not
		LightsBrakeOn    bool // Is brake light active or not
		LightsReverseOn  bool // Is reverse light active or not

		Placement PlacementData // Current truck placement in the game world

		Acceleration VectorData // Represents vehicle space linear acceleration of the truck measured in meters per second^2
		Head         VectorData // Default position of the head in the cabin space
		Cabin        VectorData // Position of the cabin in the vehicle space. This is position of the joint around which the cabin rotates. This attribute might be not present if the vehicle does not have a separate cabin
		Hook         VectorData // Position of the trailer connection hook in vehicle space
	}

	Trailer struct {
		Attached  bool          // Is the trailer attached to the truck or not
		Id        string        // Id of the cargo (internal)
		Mass      float32       // Mass of the cargo in kilograms
		Wear      float32       // Current level of trailer wear/damage between 0 (min) and 1 (max)
		Placement PlacementData // Current trailer placement in the game world
	}

	Job struct {
		Income int // Reward in internal game-specific currency

		DeadlineTime  string // Absolute in-game time of end of job delivery window. Delivering the job after this time will cause it be late
		RemainingTime string // Relative remaining in-game time left before deadline

		SourceCity    string // Localized name of the source city for display purposes
		SourceCompany string // Localized name of the source company for display purposes

		DestinationCity    string // Localized name of the destination city for display purposes
		DestinationCompany string // Localized name of the destination company for display purposes
	}

	Navigation struct {
		EstimatedTime     string // Relative estimated time of arrival
		EstimatedDistance int    // Estimated distance to the destination in meters
		SpeedLimit        int    // Current value of the "Route Advisor speed limit" in km/h
	}
}

func NewClient(BaseUrl string) Client {
	return Client{BaseUrl: BaseUrl}
}

func (c Client) GetTelemetry() (t Telemetry, err error) {
	resp, err := http.Get(c.BaseUrl + "/api/ets2/telemetry")
	if err != nil {
		log.Fatalf("HTTP GET failed: %v", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Bad status %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read http response: %v", err)
		return
	}

	t, err = parseTelemetry(body)
	return
}

func parseTelemetry(jsonIn []byte) (t Telemetry, err error) {
	err = json.Unmarshal(jsonIn, &t)
	if err != nil {
		log.Fatalf("Failed to parse: %v", jsonIn)
	}
	return
}
