from flask import Flask
from flask import Response
from flask import request

app = Flask(__name__)


@app.route('/test/*', methods=('GET', 'POST'))
@app.route('/test', methods=('GET', 'POST'))
def get():
    if request.method == 'POST':
        try:
            return request.get_json()
        except Exception as e:
            print("Fail ", e)
            return str(e)
    print(request.headers.keys())
    print(request.headers)
    resp = Response("Hello World!!! the response")
    resp.headers.add_header("Response", "check")
    return resp


if __name__ == '__main__':
    context = ('server.crt', 'server.key')
    app.run(debug=True, host='0.0.0.0', port=8080, ssl_context=context)
