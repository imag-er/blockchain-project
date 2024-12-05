#!/bin/sh
docker rmi pythonpkgs
docker build -t pythonpkgs --file pythonpkgs.dockerfile  .
docker images | grep pythonpkgs