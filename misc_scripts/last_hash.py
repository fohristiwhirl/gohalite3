s = ""

with open("logs/log-0.txt") as infile:
	for line in infile:
		if "Last known hash:" in line:
			i = line.index("Last known hash: ")
			s = line[i + len("Last known hash: "):].strip()

print(s)
