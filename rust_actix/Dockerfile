FROM ubuntu

ARG APP_NAME=rust_actix

COPY ./target/release/${APP_NAME} /usr/local/bin/${APP_NAME}

#RUN apk update && apk add gcompat

RUN useradd --no-create-home server
RUN chmod 755 /usr/local/bin/${APP_NAME}
RUN chown server:server /usr/local/bin/${APP_NAME}

USER server:server

ENV APP_NAME=${APP_NAME}

ENTRYPOINT $APP_NAME

EXPOSE 8080