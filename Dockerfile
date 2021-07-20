FROM golang:1.16

RUN mkdir -p /app
WORKDIR /app

# copy the content 
COPY . .

# install dependencies
RUN go build

# execute
CMD ["./vat"]