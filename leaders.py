import requests

data = requests.get("https://api.2018.halite.io/v1/api/leaderboard").json()

data.sort(key = lambda foo : foo["mu"], reverse = True)

print()

for item in data[:30]:

	print(" {0:>3}  {1:>18} v{2:<4}  {3:.2f} +/- {4:.2f}  ({5})".format(
		item["rank"],
		item["username"],
		item["version_number"],
		item["mu"],
		item["sigma"],
		item["language"],
	))

input()
