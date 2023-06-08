import json
import random

import requests
from flask import Flask
from flask import Response
from flask import request

app = Flask(__name__)


@app.route('/manoj/*', methods=('GET', 'POST'))
@app.route('/manoj', methods=('GET', 'POST'))
def get():
    # if request.method == 'POST':
    #   #print("Received Request")
    #   data = request.get_data()
    #   if len(data) < 500:
    #     print(data)
    #   return data
    if request.method == 'POST':
        try:
            # print(jsonify(request.get_json()))
            # print(request.get_json())
            # print(request.get_data())
            return hello_world(request)
        except Exception as e:
            print("Fail ", e)
            return str(e)

    print(request.headers.keys())
    print(request.headers)
    resp = Response("Hello World!!!")
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
    return (json.dumps(test(request)), 200, headers)


@app.route('/test', methods=('GET', 'POST'))
@app.route('/test/*', methods=('GET', 'POST'))
def salesforceoauth():
    if request.method == 'GET':
        resp = requests.post(url="https://login.salesforce.com/services/oauth2/token",
                             data={'grant_type': 'authorization_code',
                                   'code': request.args['code'],
                                   'client_id': '3MVG9y7s1kgRAI8b6hp7If35rb6MSRrDIfXcnwwDWUeHXPLngSd5ho8z6liyZJA7jxUWMiyE9.YdobE9ghii7',
                                   'client_secret': '205905459C1F7F2AC19C1614E722786D5138C36D47BBB06C0F198F5A9146D10B',
                                   'redirect_uri': 'https://got.migtunnel.net/sfoauth'
                                   })
        return resp.json()
    else:
        return Response("Hello world")

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=3032)
