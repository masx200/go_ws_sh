#!/bin/bash
cp ../*.crt ./
cp ../*.key ./
chmod 777 ./*
mkdir -p -v source
mkdir -p -v build
cp -v -f ./* ./source/
makeself source build/go_ws_sh_installer.run "go_ws_sh Installer" ./install.sh
