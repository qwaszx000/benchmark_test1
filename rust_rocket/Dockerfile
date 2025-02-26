FROM ubuntu

ARG APP_NAME=rust_rocket

COPY ./target/release/${APP_NAME} /usr/local/bin/${APP_NAME}
COPY ./Rocket.toml /usr/local/bin/Rocket.toml

#RUN apk update && apk add gcompat

RUN useradd --no-create-home server
RUN chmod 755 /usr/local/bin/${APP_NAME}
RUN chown server:server /usr/local/bin/${APP_NAME}

USER server:server

ENV APP_NAME=${APP_NAME}
ENV ROCKET_CONFIG="/usr/local/bin/Rocket.toml"

ENTRYPOINT $APP_NAME

EXPOSE 8080