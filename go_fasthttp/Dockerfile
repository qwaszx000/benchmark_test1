FROM alpine

COPY ./go_fasthttp /usr/local/bin/go_fasthttp

RUN apk update && apk add gcompat

RUN adduser --no-create-home --disabled-password server
RUN chmod 755 /usr/local/bin/go_fasthttp
RUN chown server:server /usr/local/bin/go_fasthttp

USER server:server

ENTRYPOINT [ "go_fasthttp" ]
EXPOSE 8080