FROM ubuntu

COPY ./target/release/rust_ntex /usr/local/bin/rust_ntex

#RUN apk update && apk add libgcc openssl gcompat

RUN useradd --no-create-home server
RUN chmod 755 /usr/local/bin/rust_ntex
RUN chown server:server /usr/local/bin/rust_ntex

USER server:server

ENTRYPOINT [ "rust_ntex" ]
EXPOSE 8080