import QtQuick 2.5
import QtQuick.Controls 2.5
import QtQuick.Layouts 2.5
import net.asivery.AppLoad 1.0

Rectangle {
    id: root
    anchors.fill: parent
    color: "white"
    property string currentView: "nearest"
    property real latitude: 41.9029  // Default latitude (Rome, Italy)
    property real longitude: 12.4964 // Default longitude
    property real radius: 50.0       // Default search radius in Nm
    property var aircraftInfo: {}
    property var flId: ""
    property var waypoints: [] // List of waypoints
    function showPopup() {
        popupLayer.visible = true;
    }

    function hidePopup() {
        popupLayer.visible = false;
    }
    // AppLoad component for backend communications
    AppLoad {
        id: appload
        applicationID: "reRadar24"

        Component.onCompleted: {
            requestWaypoints();
        }

        onMessageReceived: (type, contents) => {
            console.log(type + " | " + contents);
            if (type === 200) {
                console.log("Init");
                requestWaypoints();

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
                            info: `${a.model} (${a.registration}) | ${a.route} (${a.operator}) | Distance: ${Math.trunc(a.distance, 4)}Nm`,
                            flightId: a.flightId
                        });
                    });
                } else if (json_contents.category === "mostTracked") {
                    mostTrackedAircraftModel.clear();
                    json_contents.aircraft.forEach(a => {
                        console.log(a);
                        mostTrackedAircraftModel.append({
                            info: `${a.callsign} | Model: ${a.model} | Route: ${a.route} (${a.flight}) | Squawk: ${a.squawk} |  Clicks: ${a.followers}`,
                            flightId: a.flightId
                        });
                    });
                }
                break;
            case 102:
                waypoints = json_contents.waypoints;
                latitude = waypoints[0].latitude;
                longitude = waypoints[0].longitude;
                requestNearestAircraft();
                console.log("wp: ", waypoints);
                break;
            case 103:
                aircraftInfo = json_contents;
                console.log("aircraftInfo: ", aircraftInfo);
                showPopup();
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

        function requestAircraftInfo(flightId) {
            flId = flightId;
            console.log(flId);
            sendMessage(4, JSON.stringify({
                flightId: flightId
            }));
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
        anchors.bottomMargin: 20 // Margin between the list and the refresh button
        color: "transparent"
        z: 1 // Keep the list container behind the refresh button

        ListView {
            id: aircraftListView
            anchors.fill: parent
            clip: true // Clip list items to stay within the container's bounds
            model: currentView === "nearest" ? nearestAircraftModel : mostTrackedAircraftModel

            delegate: Rectangle {
                width: aircraftListView.width * 0.85
                anchors.horizontalCenter: parent.horizontalCenter
                height: Math.max(60, textItem.height) // Minimum height of 60px for list items
                border.color: "black"
                color: "white"

                Text {
                    id: textItem
                    width: parent.width  // Set width to the parent Rectangle's width to ensure wrapping
                    wrapMode: Text.WordWrap // Enable word wrapping
                    anchors.centerIn: parent
                    text: info
                    font.pointSize: 24
                    color: "black"
                    verticalAlignment: Text.AlignVCenter // Align text vertically in the center
                }

                MouseArea {
                    anchors.fill: parent
                    onClicked: {
                        appload.requestAircraftInfo(flightId);
                    }
                }
            }
        }
    }

    Rectangle {
        id: popupLayer
        anchors.fill: parent
        color: "#FFFFFF" // White background for better e-ink contrast
        visible: false
        z: 999

        Rectangle {
            id: popup
            width: parent.width * 0.8
            height: parent.height * 0.8
            // z: 999
            anchors.centerIn: parent
            color: "#ffffff" // Black popup for high contrast
            radius: 0 // No rounded corners for better e-ink rendering
            border.color: "#555555" // Dark grey border
            border.width: 2
            // Close 'X' button at the top right
            Button {
                id: closeButton
                anchors.topMargin: 4
                anchors.rightMargin: 4
                rightPadding: 4
                topPadding: 4
                text: "X"
                anchors.top: parent.top
                anchors.right: parent.right
                width: 40
                height: 40
                onClicked: hidePopup()
                contentItem: Text {
                    text: parent.text
                    font.pixelSize: 25
                    color: "#000"
                }
                background: Rectangle {
                    color: rgba(255, 255, 255, 0)
                }
            }

            Button {
                id: reloadButton
                anchors.topMargin: 4
                anchors.rightMargin: 4
                anchors.top: closeButton.bottom
                text: "Reload"
                anchors.right: parent.right
                width: 40
                height: 40
                onClicked: appload.requestAircraftInfo(flId)
                contentItem: Text {
                    horizontalAlignment: Text.AlignRight
                    text: parent.text
                    font.pixelSize: 25
                    color: "#000"
                }
                background: Rectangle {
                    color: rgba(255, 255, 255, 0)
                }
            }

            Text {
                id: aircraftInfoText
                anchors.top: parent.top
                anchors.horizontalCenter: parent.horizontalCenter
                font.pixelSize: 48 // Larger font for readability
                font.bold: true
                anchors.topMargin: 2
                text: aircraftInfo.model + " - " + aircraftInfo.registration
                color: "#000000" // Black text for contrast
            }

            Image {
                id: aircraftImage
                anchors.top: aircraftInfoText.bottom
                anchors.horizontalCenter: parent.horizontalCenter
                anchors.topMargin: 2
                width: parent.width * 0.8
                source: aircraftInfo.aircraftImageUrl
                fillMode: Image.PreserveAspectFit
                smooth: true
            }
            GridLayout {
                id: contents
                columns: 2
                anchors.top: aircraftImage.bottom
                anchors.horizontalCenter: parent.horizontalCenter
                anchors.topMargin: 10
                anchors.leftMargin: 40
                anchors.rightMargin: 40
                rowSpacing: 10
                columnSpacing: 20

                Text {
                    Layout.fillWidth: true
                    text: "Speed: " + (aircraftInfo.speed) + "kt"
                    font.pixelSize: 24
                    color: "#000000"
                }

                Text {
                    Layout.fillWidth: true
                    text: "Altitude: " + (aircraftInfo.altitude) + "ft"
                    font.pixelSize: 24
                    color: "#000000"
                }

                Text {
                    Layout.fillWidth: true
                    text: "Operator: " + (aircraftInfo.operator || "N/A")
                    font.pixelSize: 24
                    color: "#000000"
                }

                Text {
                    Layout.fillWidth: true
                    text: "Route: " + (aircraftInfo.route || "N/A")
                    font.pixelSize: 24
                    color: "#000000"
                }

                Text {
                    Layout.fillWidth: true
                    visible: aircraftInfo.callsign && aircraftInfo.callsign.length > 0
                    text: "Callsign: " + (aircraftInfo.callsign || "N/A")
                    font.pixelSize: 24
                    color: "#000000"
                }

                Text {
                    Layout.fillWidth: true
                    text: "Departure: " + (aircraftInfo.departureAirport || "N/A")
                    font.pixelSize: 24
                    color: "#000000"
                }

                Text {
                    Layout.fillWidth: true
                    text: "Arrival: " + (aircraftInfo.arrivalAirport || "N/A")
                    font.pixelSize: 24
                    color: "#000000"
                }

                Text {
                    Layout.fillWidth: true
                    text: "Country: " + (aircraftInfo.country || "N/A")
                    font.pixelSize: 24
                    color: "#000000"
                }

                Text {
                    Layout.fillWidth: true
                    text: "Takeoff: " + (aircraftInfo.takeOffTime || "N/A")
                    font.pixelSize: 24
                    color: "#000000"
                }

                Text {
                    Layout.fillWidth: true
                    text: "Landing: " + (aircraftInfo.landingTime || "N/A")
                    font.pixelSize: 24
                    color: "#000000"
                }
            }
        }
    }
}
