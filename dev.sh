#!/usr/bin/env bash

export GOPATH=$GOPATH:`pwd`
gin --excludeDir ./bin/ --bin ./bin/xnote --port=3001 -i build main.go