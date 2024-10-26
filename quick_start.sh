#!/bin/bash

set -e

osCheck=`uname -m`
if [[ $osCheck =~ 'x86_64' ]];then
    architecture="amd64"
elif [[ $osCheck == 'aarch64' ]];then
    architecture="arm64"
else
    echo "Unsupported system architecture. Please refer to the official documentation"
    exit 1
fi

if [[ "$OSTYPE" =~ ^linux ]]; then
    os="linux"
elif [[ "$OSTYPE" =~ ^darwin ]]; then
    os="darwin"
else
    echo "Unsupported operating system. Please refer to the official documentation"
    exit 1
fi

VERSION=$(curl -s https://github.com/hantbk/vbackup/release/latest)
if [[ "y${VERSION}" == "y" ]];then
    echo "Failed to retrieve the latest version. Please try again later."
    exit 1
fi

function install(){
    mv $1 /usr/local/bin/vbackup_server && chmod +x /usr/local/bin/vbackup_server
    if [[ "$OSTYPE" =~ ^linux ]]; then
        mv vbackup.service /etc/systemd/system/
        systemctl enable vbackup
        systemctl daemon-reload
        systemctl start vbackup
        for b in {1..30}
        do
            sleep 3
            service_status=`systemctl status vbackup 2>&1 | grep Active`
            if [[ $service_status == *running* ]];then
                echo "vbackup service started successfully!"
                systemctl status vbackup
                break;
            else
                echo "Error starting vbackup service!"
                exit 1
            fi
        done
    elif [[ "$OSTYPE" =~ ^darwin ]]; then
        echo "Run 'vbackup_server' in the command line to start the service."
        kubackup_server
    else
        echo "Unsupported operating system. Please refer to the official documentation"
        exit 1
    fi
}

package_file_name="vbackup_server_${VERSION}_${os}_${architecture}"
# HASH_FILE_URL="https://gitee.com/kubackup/kubackup/releases/download/${VERSION}/${package_file_name}.sum"
# package_download_url="https://gitee.com/kubackup/kubackup/releases/download/${VERSION}/${package_file_name}"
# service_file_url="https://kubackup.cn/kubackup.service"
expected_hash=$(curl -s "$HASH_FILE_URL" | awk '{print $1}')

if [ -f ${package_file_name} ];then
    actual_hash=$(sha256sum "$package_file_name" | awk '{print $1}')
    if [[ "$expected_hash" == "$actual_hash" ]];then
        echo "Installation package already exists. Skipping download."
        install ${package_file_name}
        exit 0
    else
        echo "Installation package exists, but hash values do not match. Starting re-download."
        rm -f ${package_file_name}
    fi
fi

echo "Starting download of vbackup version ${VERSION} installation package"
echo "Installation package download URL: ${package_download_url}"

curl -Lk -o ${package_file_name} ${package_download_url}

if [[ "$OSTYPE" =~ ^linux ]]; then
    curl -Lk -o vbackup.service ${service_file_url}
fi

if [ ! -f ${package_file_name} ];then
	echo "Failed to download the installation package. Please try again later."
	exit 1
fi

install ${package_file_name}
