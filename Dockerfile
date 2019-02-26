FROM golang:latest

WORKDIR /home

COPY . /home

ENTRYPOINT ["go", "run", "."]
