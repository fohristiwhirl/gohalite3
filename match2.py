import json, random, subprocess

REPLAY_FOLDER = "replays_local"

bots = [
	"bot.exe",
	"other\\v15.exe"
]

# ------------------------------------------------------------------------

scores = [0,0]
positions = [0, 1]

print("{} --- {}".format(bots[0], bots[1]))

while 1:

	random.shuffle(positions)

	tmp_positions = [bots[positions[0]], bots[positions[1]]]

	args = ["__halite.exe", "-i", REPLAY_FOLDER, "--results-as-json", "--no-logs"] + tmp_positions

	output = subprocess.check_output(args).decode("ascii")

	result = json.loads(output)

	for key in result["stats"]:
		rank = result["stats"][key]["rank"]
		i = positions[int(key)]

		if rank == 1:
			scores[i] += 1

	print(scores)
