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
    --pattern=\(\.go\|\.tmpl\)$

########
# prod #
########

FROM scratch as prod

COPY --from=development /build/main /

# port number to host server on
ARG PORT

# configuration for connecting to database
ARG DATABASE_URL
ARG DATABASE_USERNAME
ARG DATABASE_PASSWORD
ARG DATABASE_NAME

ENTRYPOINT ["/main"]

