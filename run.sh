#!/usr/bin/env bash

echo "Migrate"
./aten migrate up

echo "Start server..."
./aten
