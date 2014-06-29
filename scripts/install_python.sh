#!/bin/bash

PREFIX=$1
VERSION="3.4.1"
PYTHONTARURL="https://www.python.org/ftp/python/${VERSION}/Python-${VERSION}.tar.xz"
wget -c $PYTHONTARURL
tar Jxvf Python-${VERSION}.tar.xz
cd Python-${VERSION}
./configure --prefix=$PREFIX
make
make install
cd ..
