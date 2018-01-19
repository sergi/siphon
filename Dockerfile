FROM golang:alpine AS build-env-go
RUN apk update && apk upgrade && \
      apk add --no-cache bash git openssh
# Install vendor/deps
RUN go get github.com/golang/dep/cmd/dep
RUN mkdir -p /go/src/github.com/sergi/siphon
COPY Gopkg.lock Gopkg.toml /go/src/github.com/sergi/siphon/
WORKDIR /go/src/github.com/sergi/siphon
RUN dep ensure -vendor-only
# Build actual project
COPY . /go/src/github.com/sergi/siphon/
RUN GOOS=linux go build ./cmd/siph/

FROM node:alpine AS build-env-node
RUN apk update && apk upgrade && \
      apk add --no-cache bash git openssh
RUN git clone https://github.com/sergi/siphon-ui.git app
WORKDIR /app
RUN npm install
RUN npm run build

FROM nginx:alpine
EXPOSE 80
# UDP Port
EXPOSE 1200
# WebSockets Port
EXPOSE 3000
WORKDIR /app
COPY --from=build-env-go /go/src/github.com/sergi/siphon/siph /app/
COPY --from=build-env-node /app/build /usr/share/nginx/html/
ENTRYPOINT ["./siph"]
CMD ["server"]
