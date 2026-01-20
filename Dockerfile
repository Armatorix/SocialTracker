FROM golang:1.25.5 AS go-build


WORKDIR /go/src/github.com/Armatorix/SocialTracker/be
COPY ./be/go.mod \
    ./be/go.sum \
    ./
RUN go mod download
COPY ./be ./
RUN CGO_ENABLED=0 go build -o apibin

FROM oven/bun:1.3.3  AS node-build
WORKDIR /app/
COPY ./fe/package.json \
    ./fe/bun.lockb ./

RUN bun i

COPY ./fe/.env.production .env 
COPY ./fe/public public
COPY ./fe/index.html \
    ./fe/tsconfig.json \
    ./fe/tsconfig.node.json \
    ./fe/vite.config.ts \
    ./fe/tailwind.config.cjs \
    ./fe/postcss.config.cjs \
    ./
COPY ./fe/src src


RUN bun run build

FROM golang:1.25.1-bookworm

WORKDIR /app

RUN apt-get update -y && apt-get install ca-certificates -y

COPY --from=go-build \
    /go/src/github.com/Armatorix/SocialTracker/be/apibin \
    /app/api

COPY --from=node-build \
    /app/dist \
    /app/public


CMD ["/app/api"]