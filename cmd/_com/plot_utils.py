import datetime
from typing import List

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

def has_continuity(file, info) -> bool:
    return any([
        True
        for f, info_data in info.items()
        if "resumes" in info_data and info_data["resumes"] == file
    ])

def plot_these(
        axs,
        order, # list[str]
        data, # dict
        info, # dict
        metric,
        label_last_run_only=False,
        plot_real=True,
        plot_smoothed=True,
):
    colors_used = {}
    bests = [ max([e["wr"] for e in d[metric]]) if d[metric] else 0 for d in data.values()]
    record = max(bests)

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
        should_skip_label = has_continuity(file, info) and label_last_run_only

        if plot_real:
            p = axs.plot(
                xs_secs,
                ys,
                linewidth=0.8,
                alpha=0.3,
                **kwargs)

            colors_used[file] = p[0].get_color()

        # smooth
        if plot_smoothed:
            p = axs.plot(
                xs_secs,
                smooth(ys, .95),
                'o',
                ls='-',
                ms=4,
                markevery=[-1],
                color=colors_used[file] if file in colors_used else info[file]["kwargs"]["color"],
                label=None if should_skip_label else l,
                alpha=1 if "pruned" in file else 0.5,
                linewidth=1)
            
            colors_used[file] = p[0].get_color()