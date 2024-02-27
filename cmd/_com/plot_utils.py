import datetime
from typing import List
import plot_utils

def smooth(scalars: List[float], weight: float) -> List[float]:  # Weight between 0 and 1
    last = scalars[0]  # First value in the plot (first timestep)
    smoothed = list()
    for point in scalars:
        smoothed_val = last * weight + (1 - weight) * point  # Calculate smoothed value
        smoothed.append(smoothed_val)                        # Save it
        last = smoothed_val                                  # Anchor the last smoothed value
    return smoothed

def get_offset(file, op, info, data):
    total_offset = 0
    if "resumes" in info[file]:
        resumes = info[file]["resumes"]
        offset = data[resumes][op][-1]["delta"]
        if "at" in info[file]:
            offset *= info[file]["at"]
        total_offset += offset
        total_offset += get_offset(resumes, op, info, data)
    return total_offset

def plot_these(
        axs,
        order, # list[str]
        data, # dict
        info, # dict
        metric,
):
    colors_used = {}
    record = max([ max([e["wr"] for e in d[metric]]) for d in data.values()])

    for file in order:
        d = data[file]
        xs = [t["delta"] for t in d[metric]]
        
        kwargs = {}

        if "resumes" in info[file]:
            offset = get_offset(file, metric, info, data)
            xs = [x + offset - xs[0] for x in xs]
            kwargs["color"] = colors_used[info[file]["resumes"]]

        if "kwargs" in info[file]: kwargs = {**kwargs, **info[file]["kwargs"]}

        xs_secs = [datetime.timedelta(seconds=x).total_seconds() for x in xs]
        ys = [t["wr"] for t in d[metric]]


        # label
        m = max(ys)
        l = f"{info[file]['label']} ({round(m*100,2)})"
        if m == record: l = "$\\bf{" + l + "}$"
        kwargs["label"] = l

        p = axs.plot(
            xs_secs,
            ys,
            linewidth=0.8,
            alpha=0.3,
            **kwargs)

        colors_used[file] = p[0].get_color()

        # smooth
        axs.plot(
            xs_secs,
            plot_utils.smooth(ys, .95),
            color=colors_used[file],
            alpha=1 if "pruned" in file else 0.6,
            linewidth=1)