###############
# development #
###############

FROM golang:1.14-alpine as development

# Change working directory to /build
WORKDIR /build

# required to fetch CompileDaemon
RUN apk add git

# Hot reloading for development
RUN go get github.com/githubnemo/CompileDaemon

# Copy and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy and build source code
COPY main.go .
COPY app ./app
RUN go build -o main .

ENTRYPOINT CompileDaemon \
    --build="go build -o main" \
    --command=./main \
    --pattern=\(\.go\|\.tmpl\)$ \
    --directory=./app

########
# prod #
########

FROM scratch as prod

COPY --from=development /build/main /

ENTRYPOINT ["/main"]

