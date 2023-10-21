!/bin/bash

set -a
source /root/projects/wow-arena-leaderboard/.env
set +a

/usr/local/bin/wow-cli fetch leaderboard