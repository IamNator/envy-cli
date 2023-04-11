
# lastest golang alpine image
FROM golang:alpine

# set working directory
WORKDIR /go/src/app

# copy the current directory contents into the container at /go/src/app
COPY . .

# install dependencies
RUN go get -d -v ./...

# build the app
RUN go build -o main .

# run the app
CMD ["./main"]