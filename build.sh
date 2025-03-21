#!/bin/bash
rm -rf output
mkdir -p output/backend
cp -r icon.png manifest.json output 
rcc --binary -o output/resources.rcc application.qrc
go build .
cp reRadar24 output/backend/entry
cd ..