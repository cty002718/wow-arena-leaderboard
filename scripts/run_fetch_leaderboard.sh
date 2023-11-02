#!/bin/bash

set -a
source /home/cty/projects/wow-arena-leaderboard/.env
set +a

/usr/local/bin/wow-cli fetch leaderboard