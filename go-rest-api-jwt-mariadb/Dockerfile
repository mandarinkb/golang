FROM golang:1.16.6-alpine
RUN apk update 
# && apk add gcc git 
WORKDIR /xxx
COPY . .
RUN go build ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /xxx/app .
# สำคัญมาก ต้อง copy .env ไปอยู่ที่เดียวกับไฟล์ที่ build 
# มิฉะนั้นจะหา config ไม่เจอ
COPY .env /root/
CMD ["./app"]  
