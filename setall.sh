#!/bin/bash
set -e
chgrp -chR noordergraf .
find . -type d | xargs chmod -c g+s
chmod -cR g=u .
echo done
