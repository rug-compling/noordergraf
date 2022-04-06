#!/usr/bin/env python3

import sys

lines = []

with open('/net/noordergraf/data/prefix.ttl', 'rt') as fp:
    for line in fp:
        a = line.split()
        if len(a) != 4:
            continue
        if a[0] != '@prefix':
            continue
        if a[1] == 'll:':
            continue
        key = a[2].split('/')[2] + a[1]
        lines.append('{}\t("{}" "{}")\n'.format(key, a[1][:-1], a[2][1:-1]))

lines.sort()

sys.stdout.write('(\n')
for line in lines:
    sys.stdout.write(line.split('\t')[1])
sys.stdout.write(')\n')
