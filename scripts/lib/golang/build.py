#!/usr/bin/env python3

import subprocess
import os

def build(name, goos=None,goarch=None):
    if goos is None:
        goos = "linux"
    if goarch is None:
        goarch = "aarch64"

    exe_suffix=""
    if goos == "windows":
        exe_suffix=".exe"

    env = os.environ.copy()
    env["GOOS"]=goos
    env["GOARACH"]=goarch

    subprocess.run(["go","build","-o",f"build/{goos}-{goarch}/{name}{exe_suffix}",f"cmd/{name}/main.go"],env=env)