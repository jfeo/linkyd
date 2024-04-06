#!/bin/bash

if [[ -f /home/feo/links.json ]]; then 
    exec /home/feo/linkyd --load /home/feo/links.json
else
    exec /home/feo/linkyd
fi
