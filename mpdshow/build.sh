#!/bin/zsh

go build -o mpdshow *.go
cmx mpdshow
mv mpdshow ~/.bin
