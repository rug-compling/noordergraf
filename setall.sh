#!/bin/bash

case `pwd` in
  /net/noordergraf|/net/noordergraf/*)
    ;;
  *)
    echo Run dit script in /net/noordergraf of een subdirectory daarvan
    exit
    ;;
esac

chgrp -chR noordergraf . 2> /dev/null
find . -type d -user $USER | xargs chmod -c g+s
chmod -cR g=u . 2> /dev/null
echo done
