import requests

data = requests.get("https://api.2018.halite.io/v1/api/leaderboard").json()

for item in data:



	print("{0:>3}  {1:>18}  {2:<6}  {3:.2f}  v{4}".format(
		item["rank"],
		item["username"],
		"(" + str(item["user_id"]) + ")",
		item["mu"],
		item["version_number"],
	))

input()
