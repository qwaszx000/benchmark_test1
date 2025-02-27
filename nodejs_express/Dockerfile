FROM node:current-alpine

RUN adduser --disabled-password --home /home/server/ server

COPY ./index.js /home/server/index.js
COPY ./package.json /home/server/package.json
RUN chown -R server:server /home/server/
RUN chmod -R 700 /home/server/

WORKDIR /home/server/
USER server:server

RUN npm install

ENTRYPOINT npm start

EXPOSE 8080