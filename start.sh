#!/bin/bash

docker build -t benzinga .

docker run -p 8080:8080 benzinga
