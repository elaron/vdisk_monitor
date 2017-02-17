#!/bin/sh

rm -rf ./agent

go build . 

./agent
