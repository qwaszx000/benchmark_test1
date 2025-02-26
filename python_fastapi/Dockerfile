FROM alpine

ARG APP_NAME=go_gnet

COPY ./${APP_NAME} /usr/local/bin/${APP_NAME}

RUN apk update && apk add gcompat

RUN adduser --no-create-home --disabled-password server
RUN chmod 755 /usr/local/bin/${APP_NAME}
RUN chown server:server /usr/local/bin/${APP_NAME}

USER server:server

ENV APP_NAME=${APP_NAME}

ENTRYPOINT $APP_NAME

EXPOSE 8080