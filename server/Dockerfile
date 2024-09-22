FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod .          
COPY go.sum .          
RUN go mod tidy        
COPY . .            

RUN go build -o server .  

CMD ["./server"]