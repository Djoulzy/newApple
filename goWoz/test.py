#!/usr/bin/env python3

import wozardry
from pprint import pprint

wz = wozardry.WozDiskImage("../imgTest/Choplifter.woz")
pprint(wz.to_json())
# wz.test()

# for i in tr.nibble():
#     pprint(i)

tr = wz.seek(34)
for x in range(10):
    print(next(tr.nibble()))

tr = wz.seek(0)
for x in range(10):
    print(next(tr.nibble()))

tr = wz.seek(16)
for x in range(10):
    print(next(tr.nibble()))