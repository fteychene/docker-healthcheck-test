FROM node:6

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY package.json /usr/src/app
RUN npm install

EXPOSE 3333

HEALTHCHECK --interval=10s --timeout=3s \
  CMD curl -f http://localhost:3333/ || exit 1

COPY . /usr/src/app

CMD ["npm", "start"]
