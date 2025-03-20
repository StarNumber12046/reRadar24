import QtQuick 2.5
import QtQuick.Controls 2.5
import net.asivery.AppLoad 1.0

Rectangle {
    id: root
    anchors.fill: parent
    color: "white"
    property string currentView: "nearest"
    property real latitude: 41.9029  // Default latitude (Rome, Italy)
    property real longitude: 12.4964 // Default longitude
    property real radius: 50.0       // Default search radius in Nm

    property var waypoints: [] // List of waypoints

    // AppLoad component for backend communications
    AppLoad {
        id: appload
        applicationID: "reRadar24"

        Component.onCompleted: {
            requestWaypoints();
            requestNearestAircraft();
        }

        onMessageReceived: (type, contents) => {
            console.log(type + " | " + contents);
            if (type === 200) {
                console.log("Init");
                requestWaypoints();
                requestNearestAircraft();
                return;
            }
            var json_contents;
            json_contents = JSON.parse(contents);
            switch (type) {
            case 101:
                if (json_contents.category === "nearest") {
                    nearestAircraftModel.clear();
                    json_contents.aircraft.forEach(a => {
                        nearestAircraftModel.append({
                            info: `${a.model} (${a.registration}) | ${a.route} (${a.operator}) | Distance: ${Math.trunc(a.distance, 4)}Nm`
                        });
                    });
                } else if (json_contents.category === "mostTracked") {
                    mostTrackedAircraftModel.clear();
                    json_contents.aircraft.forEach(a => {
                        console.log(a);
                        mostTrackedAircraftModel.append({
                            info: `${a.callsign} | Model: ${a.model} | Route: ${a.route} (${a.flight}) | Squawk: ${a.squawk}`
                        });
                    });
                }
                break;
            case 102:
                waypoints = json_contents.waypoints;
                console.log("wp: ", waypoints);
                break;
            }
        }

        function requestNearestAircraft() {
            sendMessage(1, JSON.stringify({
                latitude: latitude,
                longitude: longitude,
                radius: radius
            }));
        }

        function requestMostTrackedAircraft() {
            sendMessage(2, "{}");
        }

        function requestWaypoints() {
            sendMessage(3, "{}");
        }
    }

    ListModel {
        id: nearestAircraftModel
    }

    ListModel {
        id: mostTrackedAircraftModel
    }

    // Header
    Text {
        id: header
        anchors.top: parent.top
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.topMargin: 20
        text: currentView === "nearest" ? "Nearest Aircraft" : "Most Tracked Aircraft"
        font.pointSize: 36
        color: "black"
    }

    // Navigation Menu
    Row {
        id: navMenu
        anchors.top: header.bottom
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.topMargin: 20
        spacing: 40

        Rectangle {
            width: 200
            height: 80
            border.color: "black"
            color: currentView === "nearest" ? "#CCCCCC" : "white"

            Text {
                anchors.centerIn: parent
                text: "Nearest"
                font.pointSize: 28
                color: "black"
            }

            MouseArea {
                anchors.fill: parent
                onClicked: {
                    currentView = "nearest";
                    appload.requestNearestAircraft();
                }
            }
        }

        Rectangle {
            width: 200
            height: 80
            border.color: "black"
            color: currentView === "mostTracked" ? "#CCCCCC" : "white"

            Text {
                anchors.centerIn: parent
                text: "Popular"
                font.pointSize: 28
                color: "black"
            }

            MouseArea {
                anchors.fill: parent
                onClicked: {
                    currentView = "mostTracked";
                    appload.requestMostTrackedAircraft();
                }
            }
        }
    }

    // Waypoint Selector
    ComboBox {
        id: waypointSelector
        anchors.top: navMenu.bottom
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.topMargin: 20
        model: waypoints
        textRole: "name"

        onActivated: index => {
            latitude = waypoints[index].latitude;
            longitude = waypoints[index].longitude;
            appload.requestNearestAircraft();
        }
    }

    Row {
        id: radiusInputContainer
        anchors.top: waypointSelector.bottom
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.topMargin: 20
        height: 50
        TextField {
            id: radiusInput
            width: 200
            placeholderText: "Radius (Nm)"
            text: radius

            onEditingFinished: {
                radius = parseFloat(text) || 50.0;

                radiusInput.focus = false;
                appload.requestNearestAircraft();
            }
        }
    }

    // Refresh Button - moved to define before listContainer for proper anchoring
    Rectangle {
        id: refreshButton
        anchors.bottom: parent.bottom
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.bottomMargin: 20
        width: 300
        height: 100
        border.color: "black"
        color: "white"
        z: 2 // Make sure refresh button is above list

        Text {
            anchors.centerIn: parent
            text: "Refresh"
            font.pointSize: 28
            color: "black"
        }

        MouseArea {
            anchors.fill: parent
            onClicked: {
                if (currentView === "nearest") {
                    appload.requestNearestAircraft();
                } else {
                    appload.requestMostTrackedAircraft();
                }
            }
        }
    }

    // Aircraft List View with Margin
    Rectangle {
        id: listContainer
        anchors.top: radiusInputContainer.bottom
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.bottom: refreshButton.top
        anchors.topMargin: 20
        anchors.leftMargin: 10
        anchors.rightMargin: 10
        anchors.bottomMargin: 20 // Add margin between list and refresh button
        color: "transparent"
        z: 1 // Make sure listContainer is behind refresh button

        ListView {
            id: aircraftListView
            anchors.fill: parent
            clip: true // Enable clipping to keep list items within the container
            model: currentView === "nearest" ? nearestAircraftModel : mostTrackedAircraftModel

            delegate: Rectangle {
                width: aircraftListView.width
                height: 70
                border.color: "black"
                color: "white"

                Text {
                    anchors.centerIn: parent
                    text: info
                    font.pointSize: 24
                    color: "black"
                }
            }
        }
    }
}
