#!/bin/bash
( nohup ./skillrep-server 2>&1 & echo $! > skillrep-server.pid ) | tee -a server.log &

