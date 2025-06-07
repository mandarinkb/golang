FROM golang:1.20-buster AS build
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /bin/app

FROM debian:buster-slim
COPY --from=build /bin/app /bin
EXPOSE 3000

# install cron
RUN apt-get update \
    && apt-get -y install cron 
# Copy the cron file to the container
COPY cron.d/root /etc/cron.d/root
COPY script/* /app/script/
COPY log/* /app/log/
COPY job/* /app/job/
# Set permission
RUN chmod 755 /etc/cron.d/* \
    && chmod 755 /app/script/* \
    && chmod 755 /app/log/* \
    && chmod 755 /app/job/*
# Apply the cron job
RUN crontab /etc/cron.d/root
CMD service cron start && /bin/app