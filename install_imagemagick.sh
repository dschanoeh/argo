#!/bin/sh

wget https://www.imagemagick.org/download/ImageMagick.tar.gz
mkdir imagemagick
tar xvzf ImageMagick.tar.gz -C imagemagick --strip-components 1
cd imagemagick
./configure
make
sudo make install
sudo ldconfig /usr/local/lib
