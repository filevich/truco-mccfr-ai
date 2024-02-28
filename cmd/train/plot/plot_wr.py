import argparse
import json
import matplotlib.pyplot as plt
import sys; sys.path.append('cmd/_com')
import plot_utils

plt.rcParams['savefig.dpi'] = 224
plt.rcParams['grid.color'] = 'gainsboro'

parser = argparse.ArgumentParser(description='Plot cfr train')
parser.add_argument('-i', '--input', type=str, default='/tmp/train/result.json',
                    required=False, help='.json input file')
args = parser.parse_args()

# example expected data structure
data = {
    "train_esvmccfr_a2_2p.3280483.out": {
        "random": {"wr": 1, "u": 2, "l": 3, "di": 4, "t": -1},
        "simple": {"wr": 5, "u": 6, "l": 7, "di": 8, "t": -2},
    },
}

# plot data info
info_base = {

    # 
    # train
    # 

    "train_esvmccfr_a2_2p.3280483.out": {
        "label": "esv-a2",
        "kwargs": {
            "color": "mediumpurple",
        }
    },
    "train_esvmccfr_a1_2p.3280505.out": {
        "label": "esv-a1",
        "kwargs": {
            "color": "wheat",
        }
    },
    "train_esvmccfr_a3_2p.3280535.out": {
        "label": "esv-a3",
        "kwargs": {
            "color": "cornflowerblue",
        }
    },
    "train_esvmccfr_null_2p.3280538.out": {
        "label": "esv-null",
        "kwargs": {
            "color": "darkseagreen",
        }
    },
    "train_eslmccfr_null_2p.3282695.out": {
        "label": "esl-null",
        "kwargs": {
            "color": "lightcoral",
        }
    },
    "train_eslmccfr_null_2p.3283325.out": { # <- rename
        "label": "esl-null-4t",
        "kwargs": {
            "color": "royalblue",
        }
    },



    # 
    # resume
    # 

    "resume_esvmccfr_a3_2p_2t.3293685.out": {
        "label": "esv-a3",
        "resumes": "train_esvmccfr_a3_2p.3280535.out",
        "kwargs": {
            "color": "royalblue",
        }
    },
    "resume_eslmccfr_null_2p_2t.3294059.out": {
        "label": "esl-null",
        "resumes": "train_eslmccfr_null_2p.3282695.out",
        "kwargs": {
            "color": "lightcoral",
        }
    },
    "resume_esvmccfr_null_2p_2t.3293687.out": {
        "label": "esv-null",
        "resumes": "train_esvmccfr_null_2p.3280538.out",
        "kwargs": {
            "color": "darkseagreen",
        }
    },

    # 
    # pruned
    # 

    "pruned_esvmccfr_a1_2p_1t.3325490.out": {
        "resumes": "train_esvmccfr_a1_2p.3280505.out",
        "label": "p-esv-a1",
        "kwargs": {
            "color": "gold",
        }
    },
    
    "pruned_esvmccfr_a2_2p_1t.3325556.out": {
        "resumes": "train_esvmccfr_a2_2p.3280483.out",
        "label": "p-esv-a2",
        "kwargs": {
            "color": "rebeccapurple",
        }
    },
    
    "pruned_esvmccfr_a3_2p_1t.3325559.out": {
        # "model": "final_es-vmccfr_d10h0m_D70h0m_t385690_p0_a3_2402030916.model",
        "resumes": "train_esvmccfr_a3_2p.3280535.out",
        "at": 10/70,
        "label": "p-esv-a3-10h",
        "kwargs": {
            "color": "royalblue",
        }
    },
    
    "pruned_esvmccfr_a3_2p_1t.3325561.out": {
        # "model": "final_es-vmccfr_d70h0m_D70h0m_t3468734_p0_a3_2402052116.model",
        "resumes": "train_esvmccfr_a3_2p.3280535.out",
        "at": 70/70,
        "label": "p-esv-a3-70h",
        "kwargs": {
            "color": "royalblue",
        }
    },
    
    "pruned_esvmccfr_a3_2p_1t.3325562.out": {
        # "model": "final_es-vmccfr_d40h0m_D70h0m_t1812899_p0_a3_2402041516.model",
        "resumes": "train_esvmccfr_a3_2p.3280535.out",
        "at": 40/70,
        "label": "p-esv-a3-40h",
        "kwargs": {
            "color": "royalblue",
        }
    },

    "pruned_esvmccfr_a3_2p_1t.3325567.out": {
        "resumes": "resume_esvmccfr_a3_2p_2t.3293685.out",
        "label": "p-esv-a3-140h",
        "kwargs": {
            "color": "royalblue",
        }
    },

    "pruned_esv_null_2p_1t_f10.3325610.out": {
        # "model": "final_es-vmccfr_d10h0m_D70h0m_t445443_p0_null_2402030918.model",
        "resumes": "train_esvmccfr_null_2p.3280538.out",
        "at": 10/70,
        "label": "p-esv-null-10h",
        "kwargs": {
            "color": "darkgreen",
        }
    },

    "pruned_esv_null_2p_1t_f40.3325615.out": {
        # "model": "final_es-vmccfr_d40h0m_D70h0m_t1403826_p0_null_2402041519.model",
        "resumes": "train_esvmccfr_null_2p.3280538.out",
        "at": 40/70,
        "label": "p-esv-null-40h",
        "kwargs": {
            "color": "darkgreen",
        }
    },
    
    "pruned_esv_null_2p_1t_f70.3325617.out": {
        # "model": "final_es-vmccfr_d70h0m_D70h0m_t2151893_p0_null_2402052118.model",
        "resumes": "train_esvmccfr_null_2p.3280538.out",
        "label": "p-esv-null-70h",
        "kwargs": {
            "color": "darkgreen",
        }
    },

     "pruned_esl_null_2p_1t_f60.3325612.out": {
        # "model": "final_eslmccfr_d60h1m_D70h0m_t1695315_p0_null_2402061320.model",
        "resumes": "train_eslmccfr_null_2p.3282695.out",
        "at": 60/70,
        "label": "p-esl-null-60h",
        "kwargs": {
            "color": "firebrick",
        }
    },

    "pruned_esl_null_2p_1t_r60.3325616.out": {
        # "model": "resume_eslmccfr_d60h1m_D70h0m_t2633514_p0_null_2402100446.model",
        "resumes": "resume_eslmccfr_null_2p_2t.3294059.out",
        "at": 60/70,
        "label": "p-esl-null-120h",
        "kwargs": {
            "color": "firebrick",
        }
    },

    "pruned_esv_null_2p_1t_r70.3325620.out": {
        # "model": "resume_esvmccfr_d70h0m_D70h0m_t3305210_p0_null_2402101309.model",
        "resumes": "resume_esvmccfr_null_2p_2t.3293687.out",
        "label": "p-esv-null-140h",
        "kwargs": {
            "color": "darkgreen",
        }
    },

    # InfosetRondaLarge
    "train_esl_a3_2p_1t_irl.3356924.out": {
        "label": "esl-a3-irl",
         "kwargs": {
            "color": "darkorange",
        }
    },

    "train_esl_a3_2p_1t_irxl.3383571.out": {
        "label": "esl-a3-irxl",
         "kwargs": {
            "color": "darkslategrey",
        }
    },

    "train_esl_a3_2p_1t_irxxl.3384958.out": {
        "label": "esl-a3-irxxl",
         "kwargs": {
            "color": "indigo",
        }
    },

    # InfosetRondaLarge 2nd run
    "pruned_esl_a3_2p_1t_irl_f70.3373018.out": {
        "resumes": "train_esl_a3_2p_1t_irl.3356924.out",
        "label": "p-esv-a3-irl",
         "kwargs": {
            "color": "darkorange",
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

#
#
#
    
fig, axs = plt.subplots(1, 1, figsize=(12, 6))
fig.suptitle("1 thread vs 4 threads @ esl 2p null")
# show only
show_only = [
    "train_eslmccfr_null_2p.3282695.out",
    "train_eslmccfr_null_2p.3283325.out"
]
not_show = []
info = info_base
if len(show_only): info = {k:v for k,v in info.items() if k in show_only}
if len(not_show): info = {k:v for k,v in info.items() if k not in not_show}
# order
is_resume = lambda v: "resumes" in v
order = [k for k,v in info.items() if not is_resume(v)] + [k for k,v in info.items() if is_resume(v)]
# plot
plot_utils.plot_these(axs, order, data, info, metric="simple")
# legend
axs.legend(loc='center left', bbox_to_anchor=(1, 0.5), fontsize="8")
axs.grid()
# display
plt.tight_layout()
plt.show()

#
#
#
    
fig, axs = plt.subplots(1, 1, figsize=(12, 6))
fig.suptitle("ESV vs ESL @ 2p null")
# show only
show_only = [
    "pruned_esl_null_2p_1t_r60.3325616.out",
    "resume_eslmccfr_null_2p_2t.3294059.out",
    "train_eslmccfr_null_2p.3282695.out",
    "pruned_esv_null_2p_1t_r70.3325620.out",
    "resume_esvmccfr_null_2p_2t.3293687.out",
    "train_esvmccfr_null_2p.3280538.out",
]
not_show = []
info = info_base
if len(show_only): info = {k:v for k,v in info.items() if k in show_only}
if len(not_show): info = {k:v for k,v in info.items() if k not in not_show}
# order
is_resume = lambda v: "resumes" in v
order = [k for k,v in info.items() if not is_resume(v)] + [k for k,v in info.items() if is_resume(v)]
# plot
plot_utils.plot_these(axs, order, data, info, metric="simple", label_last_run_only=True)
# legend
axs.legend(loc='center left', bbox_to_anchor=(1, 0.5), fontsize="8")
axs.grid()
# display
plt.tight_layout()
plt.show()

#
#
#

fig, axs = plt.subplots(1, 1, figsize=(12, 6))
fig.suptitle("Different pruning start time comparison (prob=1%)")
# show only
show_only = [
    # esv-irb-a3
    "train_esvmccfr_a3_2p.3280535.out",
    "resume_esvmccfr_a3_2p_2t.3293685.out",
    "pruned_esvmccfr_a3_2p_1t.3325559.out",
    "pruned_esvmccfr_a3_2p_1t.3325561.out",
    "pruned_esvmccfr_a3_2p_1t.3325562.out",
    "pruned_esvmccfr_a3_2p_1t.3325567.out",

    # esv-irb-null
    "train_esvmccfr_null_2p.3280538.out",
    "resume_esvmccfr_null_2p_2t.3293687.out",
    "pruned_esv_null_2p_1t_f10.3325610.out",
    "pruned_esv_null_2p_1t_f40.3325615.out",
    "pruned_esv_null_2p_1t_f70.3325617.out",
    "pruned_esv_null_2p_1t_r70.3325620.out",

    # esv-irb-null
    "train_eslmccfr_null_2p.3282695.out",
    "resume_eslmccfr_null_2p_2t.3294059.out",
    "pruned_esl_null_2p_1t_f60.3325612.out",
    "pruned_esl_null_2p_1t_r60.3325616.out",
]
not_show = []
info = info_base
if len(show_only): info = {k:v for k,v in info.items() if k in show_only}
if len(not_show): info = {k:v for k,v in info.items() if k not in not_show}
# order
is_resume = lambda v: "resumes" in v
order = [k for k,v in info.items() if not is_resume(v)] + [k for k,v in info.items() if is_resume(v)]
# plot
plot_utils.plot_these(axs, order, data, info, metric="simple", label_last_run_only=True, plot_real=False)
# legend
axs.legend(loc='center left', bbox_to_anchor=(1, 0.5), fontsize="8")
axs.grid()
# display
plt.tight_layout()
plt.show()

#
#
#

fig, axs = plt.subplots(1, 1, figsize=(12, 6))
fig.suptitle("InfosetRondaBase vs InfosetRondaLarge @ 2p")
# show only
show_only = [
    # InfosetRondaLarge a3
    "train_esl_a3_2p_1t_irl.3356924.out",
    "pruned_esl_a3_2p_1t_irl_f70.3373018.out",
    # InfosetRondaXLarge a3
    "train_esl_a3_2p_1t_irxl.3383571.out",
    # InfosetRondaXXLarge a3
    "train_esl_a3_2p_1t_irxxl.3384958.out",

    # InfosetRondaBase a3
    "train_esvmccfr_a3_2p.3280535.out",
    "resume_esvmccfr_a3_2p_2t.3293685.out",
    "pruned_esvmccfr_a3_2p_1t.3325567.out",
]
not_show = []
info = info_base
if len(show_only): info = {k:v for k,v in info.items() if k in show_only}
if len(not_show): info = {k:v for k,v in info.items() if k not in not_show}
# order
is_resume = lambda v: "resumes" in v
order = [k for k,v in info.items() if not is_resume(v)] + [k for k,v in info.items() if is_resume(v)]
# plot
plot_utils.plot_these(axs, order, data, info, metric="simple", label_last_run_only=True)
# legend
axs.legend(loc='center left', bbox_to_anchor=(1, 0.5), fontsize="8")
axs.grid()
# display
plt.tight_layout()
plt.show()

#
#
#

fig, axs = plt.subplots(1, 1, figsize=(12, 6))
fig.suptitle("WR mccfr 2p vs Random")
# show only
show_only = [
    "train_esvmccfr_a2_2p.3280483.out",
    "train_esvmccfr_a1_2p.3280505.out",
    "train_esvmccfr_a3_2p.3280535.out",
    # "train_esvmccfr_null_2p.3280538.out",
    "train_eslmccfr_null_2p.3282695.out",
    # "train_eslmccfr_null_2p.3283325.out",
    "resume_esvmccfr_a3_2p_2t.3293685.out",
    "resume_eslmccfr_null_2p_2t.3294059.out",
    # "resume_esvmccfr_null_2p_2t.3293687.out",
    "pruned_esvmccfr_a1_2p_1t.3325490.out",
    "pruned_esvmccfr_a2_2p_1t.3325556.out",
    # "pruned_esvmccfr_a3_2p_1t.3325559.out",
    # "pruned_esvmccfr_a3_2p_1t.3325561.out",
    # "pruned_esvmccfr_a3_2p_1t.3325562.out",
    "pruned_esvmccfr_a3_2p_1t.3325567.out",
    # "pruned_esv_null_2p_1t_f10.3325610.out",
    # "pruned_esv_null_2p_1t_f40.3325615.out",
    # "pruned_esv_null_2p_1t_f70.3325617.out",
    # "pruned_esl_null_2p_1t_f60.3325612.out",
    "pruned_esl_null_2p_1t_r60.3325616.out",
    # "pruned_esv_null_2p_1t_r70.3325620.out",
    "train_esl_a3_2p_1t_irl.3356924.out",
    "train_esl_a3_2p_1t_irxl.3383571.out",
    "train_esl_a3_2p_1t_irxxl.3384958.out",
    "pruned_esl_a3_2p_1t_irl_f70.3373018.out",
]
# not show
not_show = []
info = info_base
if len(show_only): info = {k:v for k,v in info.items() if k in show_only}
if len(not_show): info = {k:v for k,v in info.items() if k not in not_show}
# order
is_resume = lambda v: "resumes" in v
order = [k for k,v in info.items() if not is_resume(v)] + [k for k,v in info.items() if is_resume(v)]
# plot
plot_utils.plot_these(axs, order, data, info, metric="random", plot_real=False, label_last_run_only=True)
# legend
axs.legend(loc='center left', bbox_to_anchor=(1, 0.5), fontsize="8")
axs.grid()
# display
plt.tight_layout()
plt.show()

#
#
#

fig, axs = plt.subplots(1, 1, figsize=(12, 6))
fig.suptitle("WR mccfr 2p vs Simple")
# show only
show_only = [
    "train_esvmccfr_a2_2p.3280483.out",
    "train_esvmccfr_a1_2p.3280505.out",
    "train_esvmccfr_a3_2p.3280535.out",
    # "train_esvmccfr_null_2p.3280538.out",
    "train_eslmccfr_null_2p.3282695.out",
    # "train_eslmccfr_null_2p.3283325.out",
    "resume_esvmccfr_a3_2p_2t.3293685.out",
    "resume_eslmccfr_null_2p_2t.3294059.out",
    # "resume_esvmccfr_null_2p_2t.3293687.out",
    "pruned_esvmccfr_a1_2p_1t.3325490.out",
    "pruned_esvmccfr_a2_2p_1t.3325556.out",
    # "pruned_esvmccfr_a3_2p_1t.3325559.out",
    # "pruned_esvmccfr_a3_2p_1t.3325561.out",
    # "pruned_esvmccfr_a3_2p_1t.3325562.out",
    "pruned_esvmccfr_a3_2p_1t.3325567.out",
    # "pruned_esv_null_2p_1t_f10.3325610.out",
    # "pruned_esv_null_2p_1t_f40.3325615.out",
    # "pruned_esv_null_2p_1t_f70.3325617.out",
    # "pruned_esl_null_2p_1t_f60.3325612.out",
    "pruned_esl_null_2p_1t_r60.3325616.out",
    # "pruned_esv_null_2p_1t_r70.3325620.out",
    "train_esl_a3_2p_1t_irl.3356924.out",
    "train_esl_a3_2p_1t_irxl.3383571.out",
    "train_esl_a3_2p_1t_irxxl.3384958.out",
    "pruned_esl_a3_2p_1t_irl_f70.3373018.out",
]
# not show
not_show = []
info = info_base
if len(show_only): info = {k:v for k,v in info.items() if k in show_only}
if len(not_show): info = {k:v for k,v in info.items() if k not in not_show}
# order
is_resume = lambda v: "resumes" in v
order = [k for k,v in info.items() if not is_resume(v)] + [k for k,v in info.items() if is_resume(v)]
# plot
plot_utils.plot_these(axs, order, data, info, metric="simple", plot_real=False, label_last_run_only=True)
# legend
axs.legend(loc='center left', bbox_to_anchor=(1, 0.5), fontsize="8")
axs.grid()
# display
plt.tight_layout()
plt.show()