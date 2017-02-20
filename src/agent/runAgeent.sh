#!/bin/sh

rm -rf ./agent

go build . 

./agent -port 3333 -id agent100 -peerId agent200 -ip 0.0.0.0
