#Thanks to
#https://docker-fastapi-projects.readthedocs.io/en/latest/uvicorn.html

FROM python:3.9.4-buster

RUN adduser --disabled-password --home /home/server/ server

COPY ./main.py /home/server/main.py
COPY ./requirements.txt /home/server/requirements.txt

RUN chown -R server:server /home/server/
RUN chmod -R 700 /home/server/

WORKDIR /home/server/
USER server:server

RUN pip install -r ./requirements.txt && pip install "uvicorn[standard]"

ENTRYPOINT python3 -m uvicorn --host 0.0.0.0 --port 8080 --workers 2 "main:app"

EXPOSE 8080