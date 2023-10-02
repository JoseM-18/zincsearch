#use the latest golang image
FROM golang:latest

#define the working directory
WORKDIR /go/src/app

#init the go mod
RUN go mod init

# Install gin for live reload
RUN go install github.com/codegangsta/gin@latest

#Install chi for routing
RUN go get -u github.com/go-chi/chi

# Copy the current directory contents into the container at /app
COPY ./app .

# Expose port 9090 to the outside world
EXPOSE 9090