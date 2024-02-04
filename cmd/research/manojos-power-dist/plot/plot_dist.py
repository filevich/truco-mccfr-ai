import json

# example expected data structure
data = {
    "123": 456,
}

# fetch the data with:
# `go run cmd/research/manojos-power-dist/main.go`

# read it
with open('/tmp/dist-envido.json', 'r') as f: data_envido = json.loads(f.read())
with open('/tmp/dist-flor.json', 'r') as f: data_flor = json.loads(f.read())
with open('/tmp/dist-power-sum.json', 'r') as f: data_power = json.loads(f.read())

import matplotlib.pyplot as plt

fig, axs = plt.subplots(1, 3, figsize=(14, 5))
fig.suptitle("Dist Uruguayan Truco")

data = data_envido
sorted_keys = sorted(data.keys(), key=lambda x: int(x))
values = [data[key] for key in sorted_keys]
axs[0].bar(sorted_keys, values, color='g')
axs[0].set_xlabel("Envido value")
axs[0].set_ylabel("Freq.")
axs[0].set_title("Envido dist")
axs[0].grid(axis='y', linestyle='--', alpha=0.7)
axs[0].set_xticks(range(len(values)), labels=sorted_keys, rotation=40)
axs[0].locator_params(axis='x', nbins=12)

data = data_flor
sorted_keys = sorted(data.keys(), key=lambda x: int(x))
values = [data[key] for key in sorted_keys]
axs[1].bar(sorted_keys, values, color='r')
axs[1].set_xlabel("Flor value")
axs[1].set_ylabel("Freq.")
axs[1].set_title("Flor dist")
axs[1].grid(axis='y', linestyle='--', alpha=0.7)
axs[1].set_xticks(range(len(values)), labels=sorted_keys, rotation=40)
axs[1].locator_params(axis='x', nbins=12)

data = data_power
sorted_keys = sorted(data.keys(), key=lambda x: int(x))
values = [data[key] for key in sorted_keys]
axs[2].bar(sorted_keys, values, color='b')
axs[2].set_xlabel('Sum of the "power" of the 3 cards in a hand')
axs[2].set_ylabel("Freq.")
axs[2].set_title('Sum of the "power" of the 3 cards in a hand dist')
axs[2].grid(axis='y', linestyle='--', alpha=0.7)
axs[2].set_xticks(range(len(values)), labels=sorted_keys, rotation=40)
axs[2].locator_params(axis='x', nbins=12)

plt.tight_layout()
plt.show()

