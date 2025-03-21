package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reRadar24/appload"
	"reRadar24/lib"

	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/client"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/common"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/flights"
	"github.com/a-finocchiaro/go-flightradar24-sdk/webrequest"
)

type MessageType uint32

const (
	NearestAircraftRequest     MessageType = 1
	MostTrackedAircraftRequest MessageType = 2
	WaypointsRequest           MessageType = 3
	AircraftInfoRequest        MessageType = 4

	NearestAircraftResponse     MessageType = 101
	MostTrackedAircraftResponse MessageType = 101
	WaypointsResponse           MessageType = 102
	AircraftInfoResponse        MessageType = 103
)

// Example implementation of AppLoadBackend
type ReRadarBackend struct{}

func (eb *ReRadarBackend) HandleMessage(replier *appload.BackendReplier, message appload.Message) {
	fmt.Printf("Received message type: %d, contents: %s\n", message.MsgType, message.Contents)

	if message.MsgType == appload.MsgSystemTerminate {
		fmt.Println("Received termination message")
		os.Exit(0)
		return
	}

	if message.MsgType > 1000 {
		replier.SendMessage(200, "Init")
	}

	if message.MsgType == uint32(NearestAircraftRequest) {
		var body lib.NearestAircraftRequestBody
		if err := json.Unmarshal([]byte(message.Contents), &body); err != nil {
			replier.SendMessage(uint32(NearestAircraftResponse), "Error parsing request body: "+err.Error())
			return
		}

		var requester common.Requester = webrequest.SendRequest
		zone := client.GetBoundsByPoint(body.Latitude, body.Longitude, body.Radius*1852)
		aircraftNear, _ := client.GetFlightsInZone(requester, zone)

		plen := lib.FormatFeedFlight(aircraftNear.Flights, body.Latitude, body.Longitude)
		nearestAircraft, err := json.Marshal(plen)
		if err != nil {
			replier.SendMessage(uint32(NearestAircraftResponse), "Error formatting response: "+err.Error())
			return
		}

		replier.SendMessage(uint32(NearestAircraftResponse), string(nearestAircraft))
		return
	}

	if message.MsgType == uint32(MostTrackedAircraftRequest) {
		var requester common.Requester = webrequest.SendRequest
		mostTracked, _ := client.GetFR24MostTracked(requester)

		var properAircraft []flights.Fr24MostTrackedData
		for _, flight := range mostTracked.Data {
			properAircraft = append(properAircraft, flight)
		}

		mostTrackedAircraft := lib.FormatMostTracked(properAircraft)

		mostTrackedAircraftReplyJson, err := json.Marshal(mostTrackedAircraft)
		if err != nil {
			replier.SendMessage(uint32(MostTrackedAircraftResponse), "Error formatting response: "+err.Error())
			return
		}

		replier.SendMessage(uint32(MostTrackedAircraftResponse), string(mostTrackedAircraftReplyJson))
		return
	}

	if message.MsgType == uint32(WaypointsRequest) {
		waypoints := lib.GetUserWaypoints()
		waypointsReply := lib.WaypointsResponse{
			Waypoints: waypoints,
		}

		waypointsReplyJson, err := json.Marshal(waypointsReply)
		if err != nil {
			replier.SendMessage(uint32(WaypointsResponse), "Error formatting response: "+err.Error())
			return
		}

		replier.SendMessage(uint32(WaypointsResponse), string(waypointsReplyJson))
		return
	}

	if message.MsgType == uint32(AircraftInfoRequest) {
		var body lib.AircraftInfoRequestBody
		if err := json.Unmarshal([]byte(message.Contents), &body); err != nil {
			replier.SendMessage(uint32(AircraftInfoResponse), "Error parsing request body: "+err.Error())
			return
		}
		fmt.Println(body.FlightId)

		info := lib.GetAircraftInfo(body.FlightId)
		var info_json []byte
		info_json, err := json.Marshal(info)
		if err != nil {
			replier.SendMessage(uint32(AircraftInfoResponse), "Error formatting response: "+err.Error())
			return
		}
		replier.SendMessage(uint32(AircraftInfoResponse), string(info_json))
		return
	}
}

func main() {
	backend := &ReRadarBackend{}

	app, err := appload.NewAppLoad(backend)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating app: %v\n", err)
		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running app: %v\n", err)
		os.Exit(1)
	}
}
