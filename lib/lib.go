package lib

import (
	"embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/client"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/common"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/flights"
	"github.com/a-finocchiaro/go-flightradar24-sdk/webrequest"
)

//go:embed datasets/*
var datasets embed.FS
var AircraftTypes Aircraft
var Airlines []Airline

func GetAircraftTypes() Aircraft {
	if AircraftTypes != nil {
		fmt.Println("AircraftTypes already loaded")
		return AircraftTypes
	}

	byteValue, err := datasets.ReadFile("datasets/AircraftTypes.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Attempt to unmarshal the JSON
	err = json.Unmarshal(byteValue, &AircraftTypes)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		os.Exit(1)
	}

	// Optionally, print the first entry to see if it's loaded
	fmt.Printf("Loaded %d aircraft types\n", len(AircraftTypes))
	if len(AircraftTypes) > 0 {
		fmt.Println("First Aircraft:", AircraftTypes[0].ModelFullName)
	}
	return AircraftTypes
}

// LoadAirlines loads airline data from a CSV file
func LoadAirlines() ([]Airline, error) {
	if len(Airlines) > 0 {
		fmt.Println("Airlines already loaded")
		return Airlines, nil
	}

	file, err := datasets.Open("datasets/airlines.csv")
	reader := csv.NewReader(file)
	reader.Comma = ',' // Use '^' as delimiter

	// Read all records from the file
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read CSV: %v", err)
	}
	fmt.Printf("Number of records: %d\n", len(records))
	// Iterate over the records and store them in the airlines slice
	for _, record := range records {

		if len(record) != 6 {
			fmt.Printf("Skipping record with invalid column count: %v (%d)\n", record, len(record))
			continue // Skip records that don't have the expected number of columns
		}

		airline := Airline{
			Name:     record[0],
			IATA:     record[1],
			ICAO:     record[2],
			Callsign: record[3],
			Country:  record[4],
			Active:   record[5] == "Y",
		}
		Airlines = append(Airlines, airline)
	}
	fmt.Printf("Number of airlines loaded: %d\n", len(Airlines))
	return Airlines, nil
}

// toRadians converts degrees to radians.
func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 3440.069

	lat1Rad := toRadians(lat1)
	lon1Rad := toRadians(lon1)
	lat2Rad := toRadians(lat2)
	lon2Rad := toRadians(lon2)

	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	a := math.Pow(math.Sin(deltaLat/2), 2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return distance
}

func stringOr(s string, def string) string {
	if len(s) > 0 {
		return s
	}
	return def
}

func FormatFeedFlight(aircraft map[string]flights.FeedFlightData, lat, lon float64) BackendAircraftResponse {
	// Declare plane and initialize it to avoid nil values
	types := GetAircraftTypes()
	airlines, err := LoadAirlines()
	if err != nil {
		fmt.Println("Error loading airlines:", err)
		return BackendAircraftResponse{
			Category: "nearest",
			Aircraft: []BackendAircraft{},
		}
	}
	var formatted_aircraft []BackendAircraft
	for flightId, plane := range aircraft {
		// Declare plane and initialize it to avoid nil values
		var current_plane AircraftElement
		var current_airline Airline
		for _, t := range types {
			if t.Designator == plane.Aircraft_code {
				current_plane = t
				break
			}
		}
		for _, a := range airlines {
			if len(a.ICAO) == 0 {
				continue
			}
			if a.ICAO == plane.Airline_icao {
				current_airline = a
			}
		}
		// Find the aircraft type by matching the designator
		for _, t := range types {
			if t.Designator == plane.Aircraft_code {
				current_plane = t
				break
			}
		}
		formatted_plane := BackendAircraft{
			Model:        current_plane.ModelFullName,
			Route:        stringOr(plane.Origin_airport_iata, "N/A") + "->" + stringOr(plane.Destination_airport_iata, "N/A"),
			Operator:     stringOr(current_airline.Name, "private owner"),
			Registration: stringOr(plane.Registration, "N/A"),
			FlightId:     flightId,
			Distance:     Haversine(lat, lon, float64(plane.Lat), float64(plane.Long)),
		}
		formatted_aircraft = append(formatted_aircraft, formatted_plane)
	}
	return BackendAircraftResponse{
		Category: "nearest",
		Aircraft: formatted_aircraft,
	}

}

func FormatMostTracked(aircraft []flights.Fr24MostTrackedData) BackendMostTrackedAircraftResponse {
	var formatted_aircraft []BackendMostTrackedAircraft
	types := GetAircraftTypes()
	for _, plane := range aircraft {
		// Declare plane and initialize it to avoid nil values
		var current_plane AircraftElement
		// Find the aircraft type by matching the designator
		for _, t := range types {
			if t.Designator == plane.Model {
				current_plane = t
				break
			}
		}
		fmt.Println(plane.Model)
		formatted_plane := BackendMostTrackedAircraft{

			Model:    stringOr(current_plane.ModelFullName, stringOr(plane.Model, stringOr(plane.Aircraft_type, "N/A"))),
			Route:    stringOr(plane.From_iata, "N/A") + "->" + stringOr(plane.To_iata, "N/A"),
			Flight:   stringOr(plane.Flight, "private owner"),
			Squawk:   stringOr(plane.Squawk, "N/A"),
			Callsign: stringOr(plane.Callsign, "N/A"),
			FlightId: plane.Flight_id,
		}
		formatted_aircraft = append(formatted_aircraft, formatted_plane)
	}
	return BackendMostTrackedAircraftResponse{
		Category: "mostTracked",
		Aircraft: formatted_aircraft,
	}
}

func GetUserWaypoints() []Waypoint {
	var waypoints []Waypoint
	// Get the user's waypoints from a json file in the user's .config
	// directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user config directory:", err)
		return waypoints
	}
	filePath := filepath.Join(configDir, "waypoints.json")
	fmt.Println("Loading waypoints from:", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return waypoints
	}
	defer file.Close()
	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return waypoints
	}

	err = json.Unmarshal(byteValue, &waypoints)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return waypoints
	}
	return waypoints
}

func GetAircraftInfo(flightId string) AircraftInfoResponse {
	var requester common.Requester = webrequest.SendRequest
	flight, err := client.GetFlightDetails(requester, flightId)
	fmt.Println("---------GOT FLIGHT----------")
	var f, _ = json.Marshal(flight)
	fmt.Println(string(f))
	if err != nil {
		fmt.Println(err)
		return AircraftInfoResponse{}
	}
	fmt.Println("Status:" + flight.Status.Text)
	var imageUrl string
	if len(flight.Aircraft.Images.Large) > 0 {
		imageUrl = flight.Aircraft.Images.Large[0].Src
		fmt.Printf("{\"Copyright\": \"%s\", \"Src\": \"%s\", \"Link\": \"%s\", \"Source\": \"%s\"}\n", flight.Aircraft.Images.Large[0].Copyright, flight.Aircraft.Images.Large[0].Src, flight.Aircraft.Images.Large[0].Link, flight.Aircraft.Images.Large[0].Source)
	} else if len(flight.Aircraft.Images.Medium) > 0 {
		imageUrl = flight.Aircraft.Images.Medium[0].Src
	} else if len(flight.Aircraft.Images.Thumbnails) > 0 {
		imageUrl = flight.Aircraft.Images.Thumbnails[0].Src
	}
	fmt.Println(imageUrl)
	return AircraftInfoResponse{
		AircraftImageUrl: imageUrl,
		Country:          flight.Aircraft.Country.Name,
		Model:            flight.Aircraft.Model.Text,
		Registration:     flight.Aircraft.Registration,
		Route:            stringOr(flight.Airport.Origin.Code.Iata, "N/A") + "->" + stringOr(flight.Airport.Destination.Code.Iata, "N/A"),
		Operator:         stringOr(flight.Airline.Name, stringOr(flight.Owner.Name, "private owner")),
		Callsign:         flight.Identification.Callsign,
		FlightId:         flight.Identification.ID,
		DepartureAirport: flight.Airport.Origin.Name,
		ArrivalAirport:   flight.Airport.Destination.Name,
	}
}
