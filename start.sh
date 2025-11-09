#!/bin/bash

export CUBE_WORKER_HOST=localhost
export CUBE_WORKER_PORT=9000

export CUBE_MANAGER_HOST=localhost
export CUBE_MANAGER_PORT=5000

go run main.go
