FROM node:latest

WORKDIR /app

COPY package*.json /app/
RUN npm install

COPY . /app/

EXPOSE 3001
CMD ["node", "api_server.js"]
