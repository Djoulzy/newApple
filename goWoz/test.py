#!/usr/bin/env python3

import wozardry
from pprint import pprint

wz = wozardry.WozDiskImage("../imgTest/MrDo.woz")
pprint(wz.to_json())
tr = wz.seek(39.75)
# tr.find(bytes.fromhex("D5 AA 96"))
# for i in tr.nibble():
#     pprint(i)
while True:
    pprint(next(tr.nibble()))
