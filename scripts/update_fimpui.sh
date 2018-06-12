#!/bin/bash
version=$1
echo "Updating to version: $version"
wget https://storage.googleapis.com/fh-repo/fimpui_${version}_armhf.deb
sudo dpkg -i fimpui_${version}_armhf.deb
