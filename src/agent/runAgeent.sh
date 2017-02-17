#!/bin/sh

rm -rf ./agent

go build . 

./agent -port 3333 -id agent100 -ip 0.0.0.0
