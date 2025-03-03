#!/bin/sh
top -b -d 1 -U $(id -u) 2>/dev/null >> /home/server/top_output.txt;

