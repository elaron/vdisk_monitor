#!/bin/sh

rm -rf ./agent

go build . 

./agent -port 4444 -id agent200 -ip 127.0.0.1
