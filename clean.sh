#!/bin/bash

set -e

if [[ "$OSTYPE" =~ ^linux ]]; then
    systemctl stop vbackup
else
    echo "For currently unsupported operating systems, please refer to the official documentation"
    exit 1
fi

rm -rf /etc/systemd/system/vbackup.service
rm -rf /usr/local/bin/vbackup_server
rm -rf ~/.vbackup
