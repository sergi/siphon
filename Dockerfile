FROM golang:alpine AS build-env-go
RUN apk update && apk upgrade && \
      apk add --no-cache bash git openssh
WORKDIR /app
ADD . /app
RUN cd /app \
      && go get -d -v ./...\
      && cd cmd/siph \
      && GOOS=linux go build
ENTRYPOINT ./cmd/siph/siph

FROM node:alpine AS build-env-node
RUN apk update && apk upgrade && \
      apk add --no-cache bash git openssh
WORKDIR /app
RUN git clone https://github.com/sergi/siphon-ui.git
RUN cd /app/siphon-ui && npm install && npm run build


FROM nginx:alpine
# WORKDIR /app
# EXPOSE 8080
COPY --from=build-env-go /app/cmd/siph/siph /app/
COPY --from=build-env-node /app/siphon-ui/build /usr/share/nginx/html
ENTRYPOINT ["./app/siph", "server"]

