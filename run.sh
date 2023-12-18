#!/usr/bin/env bash

echo "Migrate"
cd migrator && yarn migrate
echo "Start server..."
./aten
