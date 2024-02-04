import datetime

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
    }
}

# fetch the data with:
# `rsync -avz -e 'ssh -p 10022' 'juan.filevich@cluster.uy:~/batches/out/train_*.out' /tmp/train`

# parse it with:
# `python cmd/train/plot/parse_wr.py -d /tmp/train`

# read it
with open('/tmp/train/result.json', 'r') as f:
    import json
    data = json.loads(f.read())

# show only
show_only = [
    # "train_esvmccfr_a2_2p.3280483.out",
]

if len(show_only): data = {k:v for k,v in data.items() if k in show_only}

import matplotlib.pyplot as plt

# wr

fig, axs = plt.subplots(1, 2, figsize=(14, 5))
fig.suptitle("Train 2p 72hs")

axs[0].set_title("WR vs Random bot")
for file,d in data.items():
    xs = [t["delta"] for t in d["ale"]]
    xs_hours = [str(datetime.timedelta(seconds=x)) for x in xs]
    ys = [t["wr"] for t in d["ale"]]
    l = info[file]['label']
    axs[0].plot(range(len(ys)), ys, label=l)
    axs[0].set_xticks(range(len(ys)), labels=xs_hours, rotation=40)
    axs[0].locator_params(axis='x', nbins=12)
axs[0].legend()

axs[1].set_title("WR vs Simple bot")
for file,d in data.items():
    xs = [t["delta"] for t in d["simple"]]
    xs_hours = [str(datetime.timedelta(seconds=x)) for x in xs]
    ys = [t["wr"] for t in d["simple"]]
    l = info[file]['label']
    axs[1].plot(range(len(ys)), ys, label=l)
    axs[1].set_xticks(range(len(ys)), labels=xs_hours, rotation=40)
    axs[1].locator_params(axis='x', nbins=12)
axs[1].legend()

plt.tight_layout()
plt.show()

# di

fig, axs = plt.subplots(1, 2, figsize=(14, 5))
fig.suptitle("Train 2p 72hs")

axs[0].set_title("Dumbo Index vs Random bot")
for file,d in data.items():
    xs = [t["delta"] for t in d["ale"]]
    xs_hours = [str(datetime.timedelta(seconds=x)) for x in xs]
    ys = [t["di"] for t in d["ale"]]
    l = info[file]['label']
    axs[0].plot(range(len(ys)), ys, label=l)
    axs[0].set_xticks(range(len(ys)), labels=xs_hours, rotation=40)
    axs[0].locator_params(axis='x', nbins=12)
axs[0].legend()

axs[1].set_title("Dumbo Index vs Simple bot")
for file,d in data.items():
    xs = [t["delta"] for t in d["simple"]]
    xs_hours = [str(datetime.timedelta(seconds=x)) for x in xs]
    ys = [t["di"] for t in d["simple"]]
    l = info[file]['label']
    axs[1].plot(range(len(ys)), ys, label=l)
    axs[1].set_xticks(range(len(ys)), labels=xs_hours, rotation=40)
    axs[1].locator_params(axis='x', nbins=12)
axs[1].legend()

plt.tight_layout()
plt.show()

