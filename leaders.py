import requests

COUNT = 38

data = requests.get("https://api.2018.halite.io/v1/api/leaderboard?limit={}".format(COUNT)).json()

data.sort(key = lambda foo : foo["mu"], reverse = True)

print()

for n, item in enumerate(data):

	print(" {0:>3} {1:<6}  {2:>22} v{3:<4}  {4:.2f} +/- {5:.2f}    {6}".format(
		n + 1,
		"(" + str(item["user_id"]) + ")",
		item["username"],
		item["version_number"],
		item["mu"],
		item["sigma"],
		item["language"],
	))

input()
