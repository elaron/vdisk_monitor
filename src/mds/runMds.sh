#!/bin/sh

go build -o build/runMds main.go dbConnector.go etcdInterface.go mds_model.go connector.go
./build/runMds
