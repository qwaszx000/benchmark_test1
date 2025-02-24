from flask import Flask, make_response

app = Flask(__name__)

@app.get("/test_plain")
def test_handler():
    resp = make_response("Hello world!")
    resp.content_type = "text/plain"
    return resp