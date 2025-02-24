from fastapi import FastAPI, Response

app = FastAPI()

@app.get("/test_plain")
def test_handler():
    resp = Response("Hello world!")
    resp.media_type = "text/plain"
    return resp