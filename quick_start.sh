#!/bin/bash

set -e

# Detect system architecture
osCheck=`uname -m`
if [[ $osCheck == 'x86_64' ]]; then
    architecture="amd64"
elif [[ $osCheck == 'aarch64' ]]; then
    architecture="arm64"
else
    echo "Unsupported system architecture. Please refer to the official documentation."
    exit 1
fi

# Detect OS type
if [[ "$OSTYPE" =~ ^linux ]]; then
    os="Linux"
elif [[ "$OSTYPE" =~ ^darwin ]]; then
    os="Darwin"
elif [[ "$OSTYPE" =~ ^msys ]]; then
    os="Windows"
else
    echo "Unsupported operating system. Please refer to the official documentation."
    exit 1
fi

# Retrieve the latest version tag from GitHub releases
VERSION=$(curl -s https://api.github.com/repos/hantbk/vbackup/releases/latest | grep 'tag_name' | cut -d\" -f4)
if [[ -z "$VERSION" ]]; then
    echo "Failed to retrieve the latest version. Please try again later."
    exit 1
fi

# Strip the 'v' prefix from VERSION if it exists
VERSION_STRIPPED=${VERSION#v}

# Define package file name
package_file_name="vbackup_${VERSION}_${os}_${architecture}.tar.gz"
[[ "$os" == "Windows" ]] && package_file_name="vbackup_${VERSION}_${os}_${architecture}.zip"

# Define URLs for package and checksum files
HASH_FILE_URL="https://github.com/hantbk/vbackup/releases/download/${VERSION}/vbackup_${VERSION_STRIPPED}_checksums.txt"
package_download_url="https://github.com/hantbk/vbackup/releases/download/${VERSION}/${package_file_name}"
service_file_url="https://raw.githubusercontent.com/hantbk/vbackup/main/vbackup.service"

# Fetch the expected hash from the checksum file
expected_hash=$(curl -s "$HASH_FILE_URL" | grep "$package_file_name" | awk '{print $1}')
if [[ -z "$expected_hash" ]]; then
    echo "Failed to retrieve the expected checksum. Please check the URL and try again."
    exit 1
fi

# Check if package file exists and verify checksum
if [[ -f $package_file_name ]]; then
    actual_hash=$(sha256sum "$package_file_name" | awk '{print $1}')
    if [[ "$expected_hash" == "$actual_hash" ]]; then
        echo "Installation package already exists and is verified. Skipping download."
        install $package_file_name
        exit 0
    else
        echo "Installation package exists but checksum verification failed. Re-downloading..."
        rm -f $package_file_name
    fi
fi

# Download the installation package and service file (Linux only)
echo "Downloading vbackup version ${VERSION_STRIPPED} installation package..."
curl -Lk -o ${package_file_name} ${package_download_url}

if [[ "$os" == "Linux" ]]; then
    curl -Lk -o vbackup.service ${service_file_url}
fi

# Verify if the package file downloaded successfully
if [[ ! -f $package_file_name ]]; then
    echo "Failed to download the installation package. Please try again later."
    exit 1
fi

# Re-verify the checksum after download
actual_hash=$(sha256sum "$package_file_name" | awk '{print $1}')
if [[ "$expected_hash" != "$actual_hash" ]]; then
    echo "Checksum verification failed after download. Exiting."
    exit 1
fi

# Install function
function install() {
    if [[ "$os" == "Windows" ]]; then
        unzip -o $1 -d /usr/local/bin/
    else
        tar -xzf $1 -C /usr/local/bin/ vbackup_server && chmod +x /usr/local/bin/vbackup_server
    fi

    if [[ "$os" == "Linux" ]]; then
        mv vbackup.service /etc/systemd/system/
        systemctl enable vbackup
        systemctl daemon-reload
        systemctl start vbackup
        for b in {1..30}; do
            sleep 3
            service_status=$(systemctl is-active vbackup)
            if [[ $service_status == "active" ]]; then
                echo "vbackup service started successfully!"
                systemctl status vbackup
                break
            fi
            if [[ $b -eq 30 ]]; then
                echo "Error starting vbackup service after multiple attempts!"
                exit 1
            fi
        done
    elif [[ "$os" == "Darwin" ]]; then
        echo "Run 'vbackup_server' in the command line to start the service."
        vbackup_server
    fi
}

# Proceed with installation
install $package_file_name
