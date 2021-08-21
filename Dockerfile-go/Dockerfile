FROM golang:1.16.6-alpine
# RUN apk update && apk add gcc git 
WORKDIR /xxx
COPY . .
RUN go build ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
# RUN apk --no-cache add ca-certificates
# RUN mkdir /root/data-upload

WORKDIR /root/
COPY --from=0 /xxx/app .
CMD ["./app"]  
