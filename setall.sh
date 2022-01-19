#!/bin/bash

case `pwd` in
  /net/noordergraf|/net/noordergraf/*)
    ;;
  *)
    echo Run dit script in /net/noordergraf of een subdirectory daarvan
    exit
    ;;
esac

chgrp -chR noordergraf .
find . -type d | xargs chmod -c g+s
chmod -cR g=u .
echo done
