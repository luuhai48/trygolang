#!/bin/sh

cd /app
./build migration migrate
./build run
