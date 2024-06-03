#!/usr/bin/env python3

import subprocess
import os

from mypi import repo_dir

def build(name:str):
    subprocess.run(["npm", "run", "build"],cwd=os.path.join(repo_dir,'web',name))