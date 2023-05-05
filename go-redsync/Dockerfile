FROM golang:1.18.3-alpine
RUN apk update 
# && apk add gcc git 
WORKDIR /xxx
COPY . .
RUN go build ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates &&\ 
    date && \
    apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime && \
    date
WORKDIR /root/
COPY --from=0 /xxx/app .
# สำคัญมาก ต้อง copy .env ไปอยู่ที่เดียวกับไฟล์ที่ build 
# มิฉะนั้นจะหา config ไม่เจอ
COPY config /root/config
COPY assets /root/assets
CMD ["./app"]  
