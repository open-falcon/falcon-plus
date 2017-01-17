FROM node:7

MAINTAINER masato@cepave.com

RUN npm install yarn -g

RUN useradd -m owl

USER owl

WORKDIR /home/owl

RUN git clone https://github.com/cepave-f2e/owl-light

RUN cd ./owl-light && yarn

ENV OWL_LIGHT_API_BASE=http://localhost:3000/api/v1

EXPOSE 8080

WORKDIR /home/owl/owl-light

ENTRYPOINT ["npm", "run", "dev"]
