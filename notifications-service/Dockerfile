FROM golang:1.22.0-alpine

COPY . /app

# Very important to set WORKDIR to folder where your program entry point is!
WORKDIR /app/cmd

RUN GOOS=linux go build -o main_bin main.go

CMD ["./main_bin"]
