FROM python:3.14.0a5-alpine3.21

RUN apk update && apk add gcompat

RUN adduser --disabled-password --home /home/server/ server

COPY ./main.py /home/server/main.py
COPY ./requirements.txt /home/server/requirements.txt

RUN chown -R server:server /home/server/
RUN chmod -R 700 /home/server/

WORKDIR /home/server/
USER server:server

RUN python3 -m venv venv
RUN source ./venv/bin/activate && pip install -r ./requirements.txt && pip install gunicorn

ENTRYPOINT source ./venv/bin/activate && gunicorn -w 2 -b 0.0.0.0:8080 "main:app"

EXPOSE 8080