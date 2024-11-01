#!/bin/bash
docker run -v $(pwd)/data:/data --name helper busybox true
docker cp ./data helper:/data
docker rm helper
