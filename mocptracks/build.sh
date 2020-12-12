#!/bin/zsh

go build -o mocptracks *.go
cmx mocptracks
mv mocptracks ~/.bin
