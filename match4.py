import json, random, subprocess

REPLAY_FOLDER = "replays_local"

SHUFFLE = True

bots = [
	"bot.exe",
	"other\\v15.exe",
	"other\\v15.exe",
	"node other\\dvf_29\\MyBot.js",
]

# ------------------------------------------------------------------------

scores = [0,0,0,0]
positions = [0, 1, 2, 3]

print("{} --- {} --- {} --- {}".format(bots[0], bots[1], bots[2], bots[3]))

while 1:

	if SHUFFLE:
		random.shuffle(positions)

	tmp_positions = [bots[positions[0]], bots[positions[1]], bots[positions[2]], bots[positions[3]]]

	args = ["__halite.exe", "-i", REPLAY_FOLDER, "--results-as-json", "--no-logs"] + tmp_positions

	output = subprocess.check_output(args).decode("ascii")

	result = json.loads(output)

	for key in result["stats"]:
		rank = result["stats"][key]["rank"]
		i = positions[int(key)]

		if rank == 1:
			scores[i] += 3
		elif rank == 2:
			scores[i] += 1
		elif rank == 3:
			scores[i] -= 1
		elif rank == 4:
			scores[i] -= 3

	print(scores)
