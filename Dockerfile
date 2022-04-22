FROM golang:alpine as go_builder

RUN apk add --no-cache musl-dev

WORKDIR /go/src/app
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=false go build -o app .

FROM node:alpine as asset_builder
RUN apk --no-cache add \
    --virtual .build-deps \
    	alpine-sdk \
    	cmake \
    	libssh2 libssh2-dev \
    	git \
    	dep \
    	bash \
    	curl \
    python3
WORKDIR /app
COPY package* /app/
RUN npm ci
COPY assets /app/assets
RUN npm run build

FROM alpine:3.15
WORKDIR /
RUN apk --no-cache add imagemagick
RUN addgroup -g 1000 -S app && \
    adduser -u 1000 -G app -S app
COPY --from=go_builder /go/src/app/app /
COPY templates /templates
COPY --from=asset_builder /app/assets /
USER app
ENTRYPOINT ["/app"]
