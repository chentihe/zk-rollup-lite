FROM ethereum/client-go:stable

RUN apk add --update git bash nodejs npm perl

RUN npm i -g npm@latest

RUN cd /root \
    && git clone https://github.com/cubedro/eth-net-intelligence-api.git \
    && cd eth-net-intelligence-api \
    && npm install \
    && npm install -g --unsafe-perm pm2

ADD start.sh /root/start.sh
ADD app.json /root/eth-net-intelligence-api/app.json
RUN chmod +x /root/start.sh

ENTRYPOINT /root/start.sh
