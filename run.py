import subprocess, sys

REPLAY_FOLDER = "replays_local"
IODINE_DIR = "C:\\Users\\Owner\\github\\iodine"

BOT = "bot.exe"
BOT_WITH_TESTS = "bot.exe --simtest --remaketest"

# ------------------------------------------------------------------

base_args = ["electron", IODINE_DIR, "-i", REPLAY_FOLDER]

# ------------------------------------------------------------------

def run_ref_and_quit():
	args = base_args + [BOT for n in range(4)]
	args += ["-s", "0"]
	subprocess.run(args, shell = True)
	sys.exit()

def run_tests_and_quit():
	args = base_args + [BOT_WITH_TESTS] + [BOT for n in range(3)]
	args += ["-s", "0"]
	subprocess.run(args, shell = True)
	sys.exit()

def main():

	if sys.argv[-1] == "ref":
		run_ref_and_quit()

	if sys.argv[-1] in ["test", "tests"]:
		run_tests_and_quit()

	try:
		count = int(sys.argv[-1])
	except:
		ask = input("Number of bots? ")
		if ask == "ref":
			run_ref_and_quit()
		if ask in ["test", "tests"]:
			run_tests_and_quit()
		count = int(ask)

	args = base_args + [BOT for n in range(count)]
	subprocess.run(args, shell = True)


main()
