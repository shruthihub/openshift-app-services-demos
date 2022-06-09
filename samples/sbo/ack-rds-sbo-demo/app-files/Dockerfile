FROM golang:alpine
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o aws-rds-test-api .
CMD ["/app/aws-rds-test-api"]
