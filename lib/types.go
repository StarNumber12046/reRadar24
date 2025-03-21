package lib

type Aircraft []AircraftElement

type AircraftElement struct {
	ModelFullName       string              `json:"ModelFullName"`
	Description         Description         `json:"Description"`
	Wtc                 Wtc                 `json:"WTC"`
	Wtg                 *Wtg                `json:"WTG"`
	Designator          string              `json:"Designator"`
	ManufacturerCode    string              `json:"ManufacturerCode"`
	ShowInPart3Only     bool                `json:"ShowInPart3Only"`
	AircraftDescription AircraftDescription `json:"AircraftDescription"`
	EngineCount         string              `json:"EngineCount"`
	EngineType          EngineType          `json:"EngineType"`
}

type AircraftDescription string

const (
	Amphibian  AircraftDescription = "Amphibian"
	Gyrocopter AircraftDescription = "Gyrocopter"
	Helicopter AircraftDescription = "Helicopter"
	LandPlane  AircraftDescription = "LandPlane"
	SeaPlane   AircraftDescription = "SeaPlane"
	Tiltrotor  AircraftDescription = "Tiltrotor"
)

type Description string

const (
	A1P Description = "A1P"
	A1T Description = "A1T"
	A2J Description = "A2J"
	A2P Description = "A2P"
	A2T Description = "A2T"
	A4P Description = "A4P"
	A4T Description = "A4T"
	G1P Description = "G1P"
	G1T Description = "G1T"
	H1P Description = "H1P"
	H1T Description = "H1T"
	H2P Description = "H2P"
	H2T Description = "H2T"
	H3T Description = "H3T"
	L1E Description = "L1E"
	L1J Description = "L1J"
	L1P Description = "L1P"
	L1R Description = "L1R"
	L1T Description = "L1T"
	L2E Description = "L2E"
	L2J Description = "L2J"
	L2P Description = "L2P"
	L2T Description = "L2T"
	L3J Description = "L3J"
	L3P Description = "L3P"
	L4E Description = "L4E"
	L4J Description = "L4J"
	L4P Description = "L4P"
	L4T Description = "L4T"
	L6J Description = "L6J"
	L8E Description = "L8E"
	L8J Description = "L8J"
	Lct Description = "LCT"
	S1P Description = "S1P"
	S2P Description = "S2P"
	S4P Description = "S4P"
	S4T Description = "S4T"
	T2T Description = "T2T"
	T6E Description = "T6E"
)

type EngineType string

const (
	Electric            EngineType = "Electric"
	Jet                 EngineType = "Jet"
	Piston              EngineType = "Piston"
	Rocket              EngineType = "Rocket"
	TurbopropTurboshaft EngineType = "Turboprop/Turboshaft"
)

type Wtc string

const (
	H  Wtc = "H"
	J  Wtc = "J"
	L  Wtc = "L"
	LM Wtc = "L/M"
	M  Wtc = "M"
)

type Wtg string

const (
	A Wtg = "A"
	B Wtg = "B"
	C Wtg = "C"
	D Wtg = "D"
	E Wtg = "E"
	F Wtg = "F"
	G Wtg = "G"
	Z Wtg = "Z"
)

type Airline struct {
	IATA     string `json:"IATA"`
	ICAO     string `json:"ICAO"`
	Name     string `json:"Name"`
	Callsign string `json:"Callsign"`
	Country  string `json:"Country"`
	Active   bool   `json:"Active"`
}

type BackendAircraft struct {
	Model        string  `json:"model"`
	Route        string  `json:"route"`
	Operator     string  `json:"operator"`
	Registration string  `json:"registration"`
	Distance     float64 `json:"distance"`
	FlightId     string  `json:"flightId"`
}

type BackendAircraftResponse struct {
	Category string            `json:"category"`
	Aircraft []BackendAircraft `json:"aircraft"`
}

type BackendMostTrackedAircraft struct {
	Model    string `json:"model"`
	Route    string `json:"route"`
	Flight   string `json:"flight"`
	Squawk   string `json:"squawk"`
	Callsign string `json:"callsign"`
	FlightId string `json:"flightId"`
}

type BackendMostTrackedAircraftResponse struct {
	Category string                       `json:"category"`
	Aircraft []BackendMostTrackedAircraft `json:"aircraft"`
}

type Waypoint struct {
	Name string  `json:"name"`
	Lat  float64 `json:"latitude"`
	Lon  float64 `json:"longitude"`
}

type WaypointsResponse struct {
	Waypoints []Waypoint `json:"waypoints"`
}

type NearestAircraftRequestBody struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
}

type AircraftInfoRequestBody struct {
	FlightId string `json:"flightId"`
}

type AircraftInfoResponse struct {
	AircraftImageUrl string `json:"aircraftImageUrl"`
	Country          string `json:"country"`
	Model            string `json:"model"`
	Registration     string `json:"registration"`
	Route            string `json:"route"`
	Operator         string `json:"operator"`
	Callsign         string `json:"callsign"`
	FlightId         string `json:"flightId"`
	DepartureAirport string `json:"departureAirport"`
	ArrivalAirport   string `json:"arrivalAirport"`
}
