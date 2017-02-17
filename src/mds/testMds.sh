#!/bin/sh

rm -rf ./build/testMds

go build -o build/testMds mds_model.go tdd.go

./build/testMds
