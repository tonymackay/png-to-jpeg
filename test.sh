#!/bin/bash
ARG=$1

if [ "$ARG" == "" ]; then
  ARG="."
fi

find $ARG -name '*.jpg' | xargs -n1 -I{} rm {}
find $ARG -name "*.png" -print0 | xargs -0 -I{} -P4 bash -c 'cjpeg -quality 75 -progressive -optimize -outfile "${1%.png}.jpg" $1' -- {}
