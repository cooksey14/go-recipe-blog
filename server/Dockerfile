FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod .          
COPY go.sum .          
RUN go mod tidy        
COPY . .            

RUN go build -o server .  
RUN ls -l /app
RUN ls -l /app/migrations || echo "/app/migrations does not exist"
RUN ls -l 

CMD ["./server"]