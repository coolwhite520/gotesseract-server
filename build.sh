#!/bin/sh
docker build --network=host -t ocr:1.8.6 .
docker save -o ocr.tar ocr:1.8.6