import argparse
import json
import tqdm

# update from v2.2 to v3.0

parser = argparse.ArgumentParser(description='Update .model')
parser.add_argument('-i', '--input', type=str, required=True, help='Path to input .model file')
parser.add_argument('-o', '--output', type=str, default=None, help='Path to output')
args = parser.parse_args()

def to_lower_camel_case(snake_str):
    camel_string = "".join(x.capitalize() for x in snake_str.lower().split("_"))
    return snake_str[0].lower() + camel_string[1:]

def count_lines_using_sum(file_path):
    with open(file_path, 'r') as fp:
        lines = sum(1 for _ in fp)
    return lines

def update_header(line) -> str:
    CURRENT_MODEL_VERSION = 3.0
    data = line.split(" ")
    field, values = data[0], data[1:]
    tail = " ".join(values)
    field = field.lower().replace('_', '')
    if "version" in field:
        lines = [
            f"version {CURRENT_MODEL_VERSION}",
            "hash sha160",
            "info InfosetRondaBase",
        ]
        return "\n".join(lines)
    elif "abstractor" in field: return f"abs {tail}"
    else: return f"{field.lower()} {tail}"

def update_line(line) -> str:
    return line

def check_compatibility(line):
    data = line.split(" ")
    field, values = data[0], data[1:]
    exp = "2.2"
    ok = values[0] == exp
    if not ok:
        raise Exception(f"expected version: {exp} got {values[0]}")

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
        else: print(update_line(line))
