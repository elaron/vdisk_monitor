#!/bin/sh

rm -rf ./build/runMds

go build -o build/runMds main.go mds_model.go connector.go

./build/runMds
