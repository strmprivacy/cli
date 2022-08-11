#!/usr/bin/env python3

import os
import glob
import re
import sys
import shutil

pat = re.compile(r"^(([a-z]+)_)?(([a-z]+))_?(.*)$")
os.chdir("generated_docs")
for f in glob.glob("*.md"):
    m = pat.match(f)
    if not m:
        sys.exit(1)

    _, p1, _, p2, p3 = m.groups()
    output_file = None
    path = "strm"
    if p1 and p2:
        path = os.path.join(p1,p2)
        output_file = p3
    elif p1:
        path=p1

    if p3==".md":
        output_file = "index.md"
    if not output_file :
        sys.exit(2)
    path=os.path.join(path, output_file)
    #print("%-50s %5s %10s %10s %s" % (f, p1, p2, p3, path))
    print("%-50s %s" % (f, path))
    os.makedirs(os.path.dirname(path), exist_ok=True)
    os.rename(f, path)
    if p3 == ".md":
        b = os.path.basename(os.path.dirname(path))
        try:
            shutil.copyfile(path, os.path.join(os.path.dirname(path), "..", f"{b}.md"))
        except Exception as e:
            print(e)
            pass
