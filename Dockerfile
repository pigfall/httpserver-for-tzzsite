FROM golang:1.18rc1-bullseye
COPY . /app
WORKDIR /app
RUN go build  -o httpserver .
CMD [./httpserver]
