FROM golang:1.21.5-alpine3.18 as builder

ENV CGO_ENABLED=1

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN apk update && apk upgrade
RUN apk add --no-cache sqlite build-base musl-dev nodejs npm
RUN go install github.com/a-h/templ/cmd/templ@v0.2.476
RUN npm install

RUN templ generate
RUN go build --tags 'fts5' -v -o server
RUN npm run build

CMD ["/app/server"]
