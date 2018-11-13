import os, requests, subprocess

MY_ID = 219
INITIAL_LIMIT = 100
RELOAD_LIMIT = 50
FLUORINE_EXE = "C:\\Users\\Owner\\github\\fluorine\\dist\\win-unpacked\\Fluorine.exe"

SHOW_CHALLENGES = True

class RecentGames:

	def __init__(self, my_id):

		self.game_ids = []
		self.reload(my_id, INITIAL_LIMIT)

	def reload(self, my_id, limit):

		recent = reversed(requests.get("http://api.2018.halite.io/v1/api/user/{}/match?&order_by=desc,time_played&limit={}".format(my_id, limit)).json())

		for game in recent:

			if game["challenge_id"] != None and SHOW_CHALLENGES == False:
				continue

			if game["game_id"] in self.game_ids:
				continue

			self.game_ids.append(game["game_id"])

			player_objects = game["players"]		# dict of player id --> object

			player_ranks = dict()

			player_names = []

			for pid, ob in player_objects.items():
				name = ob["username"]
				player_names.append(name)
				player_ranks[name] = ob["rank"]

			player_names.sort(key = lambda name : player_ranks[name])

			for i, name in enumerate(player_names):
				if len(name) < 16:
					player_names[i] = name + ((16 - len(name)) * " ")
				elif len(name) > 16:
					player_names[i] = name[:16]

			print("{0:>3}: {1}   ({2}x{3})   {4}".format(
				len(self.game_ids) - 1,
				game["game_id"],
				game["map_width"],
				game["map_height"],
				" ".join(player_names),
			))

	def get_game_id(self, n):
		return self.game_ids[n]

	def current_len(self):
		return len(self.game_ids)


def load_in_fluorine(filename):
	subprocess.Popen("\"{}\" \"{}\"".format(FLUORINE_EXE, filename), shell = True)

my_id = MY_ID
rg = RecentGames(my_id)

while 1:

	s = input("> ")

	if len(s) == 0:
		continue

	if s in "rR":
		rg.reload(my_id, RELOAD_LIMIT)
		continue

	if s[0] in "sS":
		try:
			my_id = int(s.split()[1])
			rg.reload(my_id, INITIAL_LIMIT)
			print("OK")
			continue
		except:
			continue

	# So, try to get a game_id to download...

	if s[0] in "dD":			# D for direct, i.e. specify exact actual id
		try:
			game_id = int(s.split()[1])
		except:
			continue
	else:
		try:
			n = int(s)
			if n < 0 or n >= rg.current_len():
				continue
			game_id = rg.get_game_id(n)
		except:
			continue

	# We got a game_id...

	if not os.path.exists("./replays/"):
		os.makedirs("./replays/")

	local_filename = "./replays/{}.hlt".format(game_id)

	if not os.path.exists(local_filename):
		url = "https://api.2018.halite.io/v1/api/user/{}/match/{}/replay".format(my_id, game_id)
		hlt = requests.get(url)
		with open(local_filename, "wb") as output:
			output.write(hlt.content)

	load_in_fluorine(os.path.abspath(local_filename))

