FROM golang:1.15.6-alpine3.12 AS build
WORKDIR /build
# Set necessary environmet variables needed for our image
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod . 
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o main .

WORKDIR /dist
RUN cp /build/main .

FROM scratch
COPY --from=build /dist/main /
ENTRYPOINT [ "/main" ]