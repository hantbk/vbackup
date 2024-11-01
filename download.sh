#!/bin/bash
set -e
version=$1
if [ -d "pkg/restic_source" ]; then
	rm -rf pkg/restic_source
fi

mkdir -p pkg/restic_source
curl -Lk https://github.com/restic/restic/archive/refs/tags/v"${version}".tar.gz -o restic.tar.gz
tar -zxvf restic.tar.gz
cp -rp restic-"${version}"/internal pkg/restic_source/rinternal
cp -rp restic-"${version}"/LICENSE pkg/restic_source/
cp -rp restic-"${version}"/VERSION pkg/restic_source/
rm -rf restic.tar.gz
rm -rf restic-"${version}"
if [[ "$OSTYPE" =~ ^linux ]]; then
	# linux
	# shellcheck disable=SC2046
	sed -i "s/\"github.com\/restic\/restic\/internal/\"github.com\/hantbk\/vbackup\/pkg\/restic_source\/rinternal/g" $(grep -rl "\"github.com\/restic\/restic\/internal" pkg/restic_source/rinternal)
elif [[ "$OSTYPE" =~ ^darwin ]]; then
	# darwin
	# shellcheck disable=SC2046
	sed -i '' "s/\"github.com\/restic\/restic\/internal/\"github.com\/hantbk\/vbackup\/pkg\/restic_source\/rinternal/g" $(grep -rl "\"github.com\/restic\/restic\/internal" pkg/restic_source/rinternal)
else
	echo "Unsupported OS: $OSTYPE"
	sed -i "s/\"github.com\/restic\/restic\/internal/\"github.com\/hantbk\/vbackup\/pkg\/restic_source\/rinternal/g" $(grep -rl "\"github.com\/restic\/restic\/internal" pkg/restic_source/rinternal)
fi
