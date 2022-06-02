#!/bin/sh

# If the web server returns something like Internal Server Error
# you can use this cgi-script to test the binary. It may be that
# it was compiled on an incompatible machine. In that case
# you will see something like:
#
#     ./index: /lib/x86_64-linux-gnu/libc.so.6: version `GLIBC_2.28' not found (required by ./index)
#
# This means, you need to compile the program on the web server.
# You can use the cgi-script 'make.sh' to do this.

echo Content-type: text/plain
echo

./index 2>&1
