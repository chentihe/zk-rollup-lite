FROM node:alpine

RUN apk add --update git

RUN git clone https://github.com/cubedro/eth-netstats

WORKDIR /eth-netstats

RUN npm install
RUN npm install -g --unsafe-perm grunt-cli
RUN grunt

EXPOSE 3000

CMD ["npm", "start"]