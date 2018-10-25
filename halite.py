
REPLAY_FOLDER = "replays_local"

EXE = "C:\\Users\\Owner\\github\\dubnium\\dubnium.exe"

# ------------------------------------------------------------------

import json, subprocess, sys, time

add = ["--no-compression", "--no-logs", "--results-as-json", "-i", REPLAY_FOLDER]

if sys.argv[0] == "python":
	new_args = sys.argv[2:] + add
else:
	new_args = sys.argv[1:] + add

for n, arg in enumerate(new_args):
	new_args[n] = '"' + arg + '"'

new_args_string = " ".join(new_args)
cmd = "{} {}".format(EXE, new_args_string)

start_time = time.time()
output = subprocess.check_output(cmd).decode("ascii")
elapsed_time = time.time() - start_time

j = json.loads(output)

useful = dict()

for key in ["map_width", "map_height", "replay", "stats", "map_seed"]:
	if key in j:
		useful[key] = j[key]

useful["time"] = time.strftime("%H:%M:%S", time.gmtime(elapsed_time))

print(json.dumps(useful, indent = 4))
