REPLAY_FOLDER = "replays_local"
EXE = "__halite.exe"

# ------------------------------------------------------------------
# The idea is we invoke this script with whatever optional flags...
# ./halite.py bot.exe bot.exe -w 64 -h 64 -s 123

import json, subprocess, sys, time

args = [EXE, "-i", REPLAY_FOLDER, "--no-logs", "--no-compression", "--results-as-json"]

if sys.argv[0] == "python":
	args += sys.argv[2:]
else:
	args += sys.argv[1:]

# Official Halite doesn't have -w and -h flags...

args = [s if s != "-w" else "--width" for s in args]
args = [s if s != "-h" else "--height" for s in args]

start_time = time.time()

raw = subprocess.check_output(args, shell=True)
parsed = json.loads(raw)

useful = dict()

useful["engine"] = EXE

for key in ["map_width", "map_height", "map_seed", "replay", "stats"]:
	useful[key] = parsed[key]

useful["time"] = time.strftime("%H:%M:%S", time.gmtime(time.time() - start_time))

print(json.dumps(useful, indent = 4))
