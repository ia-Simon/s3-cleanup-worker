FROM golang:1.17-alpine
LABEL maintainer "Towty Soluções"


WORKDIR /code

RUN touch /var/log/s3_cleanup.log && \ 
    echo "40 2 1 * * s3_cleanup_worker >> /var/log/s3_cleanup.log 2>&1" > /tmp/crontab.txt && \
    crontab /tmp/crontab.txt && \
    rm /tmp/crontab.txt && \
    crond

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

CMD crond && tail -f /var/log/s3_cleanup.log

COPY . .
RUN go install