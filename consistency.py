one = "__halite.exe"
two = "C:\\Users\\Owner\\github\\dubnium\\dubnium.exe"

bot = "bot.exe"

import json, random, subprocess, sys

random.seed()

while 1:

	game_seed = random.randint(0, 2000000000)

	args = [bot for n in range(4)] + ["--no-logs", "--no-replay", "--results-as-json", "-s", str(game_seed)]

	result_one = subprocess.check_output([one] + args, shell=True)
	result_two = subprocess.check_output([two] + args, shell=True)

	j_one = json.loads(result_one)
	j_two = json.loads(result_two)

	print("Seed {}".format(game_seed))

	fail = False

	for key in j_one["stats"]:
		if j_one["stats"][key]["score"] != j_two["stats"][key]["score"]:
			fail = True

	print("     pass" if not fail else "     FAIL <-----------------------------------")
