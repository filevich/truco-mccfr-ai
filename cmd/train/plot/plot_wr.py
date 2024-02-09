import argparse
import json
import datetime
import matplotlib.pyplot as plt

parser = argparse.ArgumentParser(description='Plot cfr train')
parser.add_argument('-i', '--input', type=str, default='/tmp/train/result.json', required=False, help='.json input file')
args = parser.parse_args()

# example expected data structure
data = {
    "train_esvmccfr_a2_2p.3280483.out": {
        "ale": {"wr": 1, "u": 2, "l": 3, "di": 4, "t": -1},
        "simple": {"wr": 5, "u": 6, "l": 7, "di": 8, "t": -2},
    },
}
# plot data info
info = {
    "train_esvmccfr_a2_2p.3280483.out": {
        "label": "esv-a2",
    },
    "train_esvmccfr_null_2p.3280538.out": {
        "label": "esv-null",
    },
    "train_esvmccfr_a1_2p.3280505.out": {
        "label": "esv-a1",
    },
    "train_esvmccfr_a3_2p.3280535.out": {
        "label": "esv-a3",
    },
    "train_eslmccfr_null_2p.3282695.out": {
        "label": "esl-null",
    },
    "train_eslmccfr_null_2p.3283325.out": {
        "label": "esl-null-4t",
    },
    # resume
    "resume_eslmccfr_null_2p_2t.3294059.out": {
        "label": "esl-null",
        "resumes": "train_eslmccfr_null_2p.3282695.out",
        "kwargs": {
        }
    },
    "resume_esvmccfr_a3_2p_2t.3293685.out": {
        "label": "esv-a3",
        "resumes": "train_esvmccfr_a3_2p.3280535.out",
        "kwargs": {
        }
    },
    "resume_esvmccfr_null_2p_2t.3293687.out": {
        "label": "esv-null",
        "resumes": "train_esvmccfr_null_2p.3280538.out",
        "kwargs": {
        }
    },
}

# fetch the data with:
# `rsync -avz -e 'ssh -p 10022' 'juan.filevich@cluster.uy:~/batches/out/train_*.out' /tmp/train`

# parse it with:
# `python cmd/train/plot/parse_wr.py -d /tmp/train`

# read it
with open(args.input, 'r') as f:
    data = json.loads(f.read())

# show only
show_only = [
    # "train_esvmccfr_a3_2p.3280535.out",
    # "resume_esvmccfr_a3_2p_2t.3293685.out",
    # "resume_esvmccfr_null_2p_2t.3293687.out",
    # "train_esvmccfr_null_2p.3280538.out"
]

# skip
not_show = [
    "train_eslmccfr_null_2p.3283325.out",
]

if len(show_only): info = {k:v for k,v in info.items() if k in show_only}
if len(not_show): info = {k:v for k,v in info.items() if k not in not_show}

# order
is_resume = lambda v: "resumes" in v
order = [k for k,v in info.items() if not is_resume(v)] + [k for k,v in info.items() if is_resume(v)]

# wr

fig, axs = plt.subplots(2, 1, figsize=(10, 8))
fig.suptitle("Train 2p 72hs")

colors_used = {}

axs[0].set_title("(a) WR vs Random bot")

for file in order:
    d = data[file]
    xs = [t["delta"] for t in d["ale"]]
    
    kwargs = {}

    if "resumes" in info[file]:
        xs = [x + data[info[file]["resumes"]]["ale"][-1]["delta"] - xs[0] for x in xs]
        kwargs["color"] = colors_used[info[file]["resumes"]]

    if "kwargs" in info[file]: kwargs = {**kwargs, **info[file]["kwargs"]}

    xs_secs = [datetime.timedelta(seconds=x).total_seconds() for x in xs]
    ys = [t["wr"] for t in d["ale"]]
    m = max(ys)
    l = f"{info[file]['label']} ({round(m*100,2)})"
    kwargs["label"] = l

    p = axs[0].plot(xs_secs, ys, linewidth=0.8, **kwargs)

    colors_used[file] = p[0].get_color()

# x axis
xs = set()
for file in order:
    d = data[file]
    ts = [t["delta"] for t in d["ale"]]
    if "resumes" in info[file]:
        ts = [t + data[info[file]["resumes"]]["ale"][-1]["delta"] - ts[0] for t in ts]
    xs = xs.union(ts)
xs = sorted(xs)
xs_hours = [str(datetime.timedelta(seconds=x)) for x in xs]
axs[0].set_xticks(xs, labels=xs_hours, rotation=40)
axs[0].locator_params(axis='x', nbins=12)

# legend
axs[0].legend(loc='center left', bbox_to_anchor=(1, 0.5), fontsize="8")

axs[1].set_title("(b) WR vs Simple bot")

for file in order:
    d = data[file]
    xs = [t["delta"] for t in d["simple"]]
    
    kwargs = {}

    if "resumes" in info[file]:
        xs = [x + data[info[file]["resumes"]]["simple"][-1]["delta"] - xs[0] for x in xs]
        kwargs["color"] = colors_used[info[file]["resumes"]]

    if "kwargs" in info[file]: kwargs = {**kwargs, **info[file]["kwargs"]}

    xs_secs = [datetime.timedelta(seconds=x).total_seconds() for x in xs]
    ys = [t["wr"] for t in d["simple"]]
    m = max(ys)
    l = f"{info[file]['label']} ({round(m*100,2)})"
    kwargs["label"] = l

    p = axs[1].plot(xs_secs, ys, linewidth=0.8, **kwargs)

    colors_used[file] = p[0].get_color()

# x axis
xs = set()
for file in order:
    d = data[file]
    ts = [t["delta"] for t in d["simple"]]
    if "resumes" in info[file]:
        ts = [t + data[info[file]["resumes"]]["simple"][-1]["delta"] - ts[0] for t in ts]
    xs = xs.union(ts)
xs = sorted(xs)
xs_hours = [str(datetime.timedelta(seconds=x)) for x in xs]
axs[1].set_xticks(xs, labels=xs_hours, rotation=40)
axs[1].locator_params(axis='x', nbins=12)

# legend
axs[1].legend(loc='center left', bbox_to_anchor=(1, 0.5), fontsize="8")

plt.tight_layout()
plt.show()

