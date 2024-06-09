FROM golang:1.22.4-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /post_comment_system

EXPOSE 8080

CMD [ "/post_comment_system" ]