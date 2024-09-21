FROM golang:1.20 AS builder

WORKDIR /app
COPY go.mod .          
COPY go.sum .          
RUN go mod tidy        
COPY . .            

RUN go build -o server .  
RUN ls -l /app

CMD ["./server"]