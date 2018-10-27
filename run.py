import subprocess, sys

REPLAY_FOLDER = "replays_local"
IODINE_DIR = "C:\\Users\\Owner\\github\\iodine"

BOT = "bot.exe"

# ------------------------------------------------------------------

args = ["electron", IODINE_DIR, "-i", REPLAY_FOLDER]

if sys.argv[-1] == "ref":
	args += [BOT for n in range(4)]
	args += ["-s", "0"]
else:
	try:
		count = int(sys.argv[-1])
	except:
		print("Need number of bots")
		sys.exit()
	args += [BOT for n in range(count)]

subprocess.run(args, shell = True)

