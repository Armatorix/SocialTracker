FROM golang:1.25.4 AS go-build

WORKDIR /go/src/github.com/Armatorix/SocialTracker/be
COPY ./be/go.mod \
    ./be/go.sum \
    ./
RUN go mod download
COPY ./be ./
RUN CGO_ENABLED=0 go build -o apibin

FROM oven/bun:1.3.3 AS bun-build
WORKDIR /app/
COPY ./fe/package.json \
    ./fe/bun.lock ./

RUN bun i

COPY ./fe/tsconfig.json \
    ./fe/build.ts \
    ./
COPY ./fe/src src

RUN bun run build

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=go-build \
    /go/src/github.com/Armatorix/SocialTracker/be/apibin \
    /app/api

COPY --from=bun-build \
    /app/dist \
    /app/public

CMD ["/app/api"]