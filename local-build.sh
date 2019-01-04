#!/bin/bash
export REDIS_HOST=localhost
export REDIS_PORT=6379
export INTERNAL_API_KEY=9239

go build -o $HOME/go/bin/robolucha-publisher
$HOME/go/bin/robolucha-publisher