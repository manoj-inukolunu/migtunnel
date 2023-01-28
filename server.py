import json
import random

from flask import Flask, jsonify
from flask import Response
from flask import request

app = Flask(__name__)


@app.route('/manoj/*', methods=('GET', 'POST'))
@app.route('/test', methods=('GET', 'POST'))
def get():
    # if request.method == 'POST':
    #   #print("Received Request")
    #   data = request.get_data()
    #   if len(data) < 500:
    #     print(data)
    #   return data
    if request.method == 'POST':
        try:
            print(jsonify(request.get_json()))
            print(request.get_json())
            # print(request.get_data())
            return hello_world(request)
        except Exception as e:
            print("Fail ", e)
            return str(e)

    print(request.headers.keys())
    print(request.headers)
    resp = Response("")
    resp.headers['x-amz-apigw-id'] = "asdf"
    return resp


@app.route('/predict', methods=('GET', 'POST'))
def predict():
    if request.method == 'POST':
        return hello_world(request)
    resp = Response("")
    resp.headers['x-amz-apigw-id'] = "asdf"
    return resp


def test(request):
    request_json = request.get_json()
    instances = list()
    preds = {"predictions": []}
    count = 0
    if isinstance(request_json, str):
        data = json.loads(request_json)
    else:
        data = request_json

    for feature in data['instances']:
        preds["predictions"].append({
            "score": random.uniform(0, 1)
        })
        count = count + 1
    print(count)
    return preds


def hello_world(request):
    headers = {'Content-Type': 'application/json'}
    # print(request.get_json())
    return (json.dumps(test(request)), 503, headers)


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=3131)
