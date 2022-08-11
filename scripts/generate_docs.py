#!/usr/bin/env python3

import os
import glob
import re
import sys
import shutil

os.makedirs("generated_docs", exist_ok=True)
os.chdir("generated_docs")

for f in glob.glob("*.md"):
    comps = f.split('_')
    path = "/".join(comps)
    folder = os.path.dirname(path)
    if folder:
        os.makedirs(folder, exist_ok=True)
        os.rename(f, path)

for f in glob.glob("**/*.md", recursive=True):
    b,_ = os.path.splitext(f)
    print(f, b)
    if os.path.isdir(b):
        os.rename(f, os.path.join(b, "index.md"))

