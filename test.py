import json

import requests

count = 0
for j in range(10):
    instances = {"instances": []}
    for i in range(1000000):
        instances["instances"].append({
            "features": [4.9281602056057885, 6, 5.673203040696216, 3, "PA", 163]
        })
    response = requests.post("https://got.migtunnel.net/test", json=json.dumps(instances))
    print(response.json())
    count = count + 1

print("Done ", count)
# print(json.dumps(instances))
