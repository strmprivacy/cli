#!/usr/bin/env python3

import os
import glob
import re
import sys
import shutil

os.makedirs("generated_docs", exist_ok=True)
os.chdir("generated_docs")

# we keep a set of directories, because all references to (<directory>.md)
# have to be replaced by <directory>/index.md
directories={'strm'} # top level directory

for f in glob.glob("*.md"):
    comps = f.split('_')
    path = "/".join(comps)
    folder = os.path.dirname(path)
    if folder:
        os.makedirs(folder, exist_ok=True)
        directories.add(folder)
        os.rename(f, path)

for f in glob.glob("**/*.md", recursive=True):
    b,_ = os.path.splitext(f)
    if os.path.isdir(b):
        os.rename(f, os.path.join(b, "index.md"))

# build the regular expression
directories = "|".join(directories)
directories = f"(/cli-reference/({directories}))(?:\.md)"
pattern = re.compile(directories)

for f in glob.glob("**/*.md", recursive=True):
    with open(f,'r') as _f:
        content=_f.read()
    with open(f,'w') as _f:
        _f.write(pattern.sub(r'\1/index.md', content))
