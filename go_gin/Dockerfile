FROM alpine

COPY ./go_gin /usr/local/bin/go_gin

RUN apk update && apk add gcompat

RUN adduser --home /home/server --disabled-password server
RUN chmod 755 /usr/local/bin/go_gin
RUN chown server:server /usr/local/bin/go_gin

USER server:server

ENTRYPOINT [ "go_gin" ]
EXPOSE 8080