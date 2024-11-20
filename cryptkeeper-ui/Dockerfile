# syntax=docker/dockerfile:1

# Build stage
FROM node:14-alpine as build

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

COPY . ./



RUN npm run build

# Production stage
FROM nginx:alpine

COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
COPY default.conf /etc/nginx/conf.d/default.conf

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
