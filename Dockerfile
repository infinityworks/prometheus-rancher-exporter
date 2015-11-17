FROM gliderlabs/alpine:3.1
MAINTAINER barwell

RUN apk-install nodejs

WORKDIR /app

ADD app.js /app/
ADD package.json /app/
RUN npm install

ENV DEBUG re
EXPOSE 9010

CMD ["npm", "start"]
