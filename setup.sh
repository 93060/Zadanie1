#!/bin/bash 
PLATFORM=$(uname -m) 

if [ "$PLATFORM" = "x86_64" ]
then export GOARCH=amd64
elif [ "$PLATFORM" = "aarch64" ]
then export GOARCH=arm64
elif [ "$PLATFORM" = "armv7l" ]
then export GOARCH=arm
fi