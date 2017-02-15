#!/bin/sh

go build -o build/testMds dbConnector.go etcdInterface.go mds_model.go tdd.go
./build/testMds
