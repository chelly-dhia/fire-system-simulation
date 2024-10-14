# Dockerfile for web service
FROM node:20

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

# Set NODE_OPTIONS to allow legacy OpenSSL behavior
ENV NODE_OPTIONS=--openssl-legacy-provider

RUN npm run build

CMD ["npm", "start"]
