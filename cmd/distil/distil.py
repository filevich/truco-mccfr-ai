import argparse
import json
import tqdm

# disitil v3 model

parser = argparse.ArgumentParser(description='Update .model')
parser.add_argument('-i', '--input', type=str, required=True, help='Path to input .model file')
args = parser.parse_args()

def update_header(line) -> str:
    CURRENT_MODEL_VERSION = "3.0-distil"
    data = line.split(" ")
    field, values = data[0], data[1:]
    tail = " ".join(values)
    field = field.lower().replace('_', '')
    if "version" in field: return f"version {CURRENT_MODEL_VERSION}"
    else: return line

def update_line(line:str) -> str:
    if not line: return line
    data = line.split(" ")
    hash, tail = data[0], data[1:]
    tail = " ".join(tail)
    rnode = json.loads(tail)
    all_the_same = len(set(rnode["strategySum"])) == 1
    if all_the_same:
        # prune this line
        return None
    m = max(rnode["strategySum"])
    ix = rnode["strategySum"].index(m)
    return f"{hash} {ix}"

def check_compatibility(line):
    data = line.split(" ")
    field, values = data[0], data[1:]
    exp = "3.0"
    ok = values[0] == exp
    if not ok:
        raise Exception(f"expected version: {exp} got {values[0]}")

if __name__ == "__main__":
    num_lines = 0
    with open(args.input) as infile:
        for line in infile: num_lines += 1

    with open(args.input) as infile:
        n = 0
        for i, line in enumerate(tqdm.tqdm(infile, total=num_lines, desc="Processing lines")):
            line = line.strip()
            if line == "": n += 1
            if i == 0: check_compatibility(line)
            if n == 0: print(update_header(line))
            else:
                new_line = update_line(line)
                if new_line == None: continue
                print(new_line)
