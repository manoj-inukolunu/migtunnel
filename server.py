from flask import Flask
from flask import request
from flask import Response
import json
import random

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
    return hello_world(request)
  resp = Response("")
  resp.headers['x-amz-apigw-id']="asdf"
  return resp

@app.route('/predict', methods=('GET', 'POST'))
def predict():
  if request.method == 'POST':
    return hello_world(request)
  resp = Response("")
  resp.headers['x-amz-apigw-id']="asdf"
  return resp


def test(request):
    request_json = request.get_json()
    instances=list()
    preds = {"predictions":[]}
    count = 0
    for feature in request_json['instances']:
        preds["predictions"].append({
          "score": random.uniform(0, 1)
        })
        count = count+1
    print(count)
    return preds

def hello_world(request):
    headers = {'Content-Type': 'application/json'}
    return (json.dumps(test(request)),200,headers)


if __name__ == '__main__':
  app.run(debug=True, host='0.0.0.0', port=3131)
