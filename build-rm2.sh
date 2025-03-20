#!/bin/bash
rm -rf output
mkdir -p output/backend
cp -r icon.png manifest.json datasets output 
rcc --binary -o output/resources.rcc application.qrc
env GOOS=linux GOARCH=arm GOARM=7 go build
cp reRadar24 output/backend/entry
cd ..