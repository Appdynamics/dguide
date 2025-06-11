#!/bin/sh

set -e


error() {
    echo "Error: $1"
    exit 1
}


OS=$(uname -s)
ARCH=$(uname -m)

BINARY_NAME="dguide"

DEST_DIR="/usr/local/bin"

# Check if the destination directory is writable
if [ ! -w "$DEST_DIR" ]; then
    if [ "$OS" = "Linux" ] || [ "$OS" = "Darwin" ]; then
        echo "The destination directory $DEST_DIR is not writable. Trying with sudo..."
        sudo mv $BINARY_NAME $DEST_DIR || error "Failed to move the binary to $DEST_DIR"
    else
        error "The destination directory $DEST_DIR is not writable. Please run the script with appropriate permissions."
    fi
else
    mv $BINARY_NAME $DEST_DIR || error "Failed to move the binary to $DEST_DIR"
fi

if command -v $BINARY_NAME &> /dev/null; then
    echo "$BINARY_NAME has been successfully installed"
else
    error "Failed to verify the installation of $BINARY_NAME"
fi

echo "Installation completed successfully."
