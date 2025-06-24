#!/bin/bash

# Default values
CONFIG_FILE="./dist/config-outside.json"
NO_RESTART=false

# Parse command line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --config) CONFIG_FILE="$2"; shift ;;
        --no-restart) NO_RESTART=true ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

# Function to run the trading bot
run_sing_box() {
    go build -trimpath -o dist/sing-box.exe -ldflags '-s -buildid= -X github.com/sagernet/sing-box/constant.Version=1.0.0' ./cmd/sing-box
    ./dist/sing-box.exe run -c "$CONFIG_FILE"
}

# Main loop to handle unexpected exits
while true; do
    run_sing_box

    if [ "$NO_RESTART" = true ]; then
        echo "Exiting..."
        break
    fi
    echo "sing-box exited unexpectedly. Restarting..."

    sleep 1
    git pull
    sleep 1
done