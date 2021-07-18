#!/bin/bash
set -e

fullfile=$1
filename=$(basename -- "$fullfile")
extension="${filename##*.}"
filename="${filename%.*}"

# A4
# 3508 x 2480 px at 300dpi
# convert $1 -resize 3508x2480 -gravity center -extent 3508x2480 -units PixelsPerInch -density 300x300 output.pdf

# A5
# 2480 x 1748 px. at 300dpi
convert $1 -resize 2480x1748 -gravity center -extent 2480x1748 -units PixelsPerInch -density 300x300 "${filename}_a5.png"
