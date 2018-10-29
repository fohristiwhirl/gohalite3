REPLAY_FOLDER = "replays_local"
IODINE_DIR = "C:\\Users\\Owner\\github\\iodine"

# ------------------------------------------------------------------
# The idea is we invoke this script with whatever optional flags...
# ./iodine.py bot.exe bot.exe -w 64 -h 64 -s 123

import subprocess, sys

args = ["electron", IODINE_DIR, "-i", REPLAY_FOLDER]

if sys.argv[0] == "python":
	args += sys.argv[2:]
else:
	args += sys.argv[1:]

subprocess.run(args, shell=True)
