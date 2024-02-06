import argparse
import json
import tqdm

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
    CURRENT_MODEL_VERSION = 2.1
    data = line.split(" ")
    field, values = data[0], data[1:]
    tail = " ".join(values)
    field = field.lower().replace('_', '')
    # version 1.0
    if "prot" in field: return f"version {CURRENT_MODEL_VERSION}"
    # version 2.0
    elif "version" in field: return f"version {CURRENT_MODEL_VERSION}"
    # version <=2.0
    elif "id" in field: return f"trainer {tail.replace('-', '')}"
    else: return f"{to_lower_camel_case(field)} {tail}"

def update_line(line) -> str:
    if line == "": return ""
    xs = line.split(" ")
    hash = xs[0]
    data = " ".join(xs[1:])
    data = json.loads(data)
    data = {to_lower_camel_case(k):v for k,v in data.items()}
    return f"{hash} {json.dumps(data)}"

num_lines = 0
with open(args.input) as infile:
    for line in infile: num_lines += 1

with open(args.input) as infile:
    n = 0
    for line in tqdm.tqdm(infile, total=num_lines, desc="Processing lines"):
        line = line.strip()
        if line == "": n += 1
        if n == 0: print(update_header(line))
        else: print(update_line(line))
