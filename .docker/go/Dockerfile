FROM golang:1.19 AS compiling_stage
RUN mkdir -p /go/src/task-tracker
WORKDIR /go/src/task-tracker
COPY ./ .
RUN go install .

FROM alpine:latest
LABEL version="1.0.0"
LABEL maintainer="Ivan Ivanov<test@test.ru>"
WORKDIR /root/
COPY --from=compiling_stage /go/bin/task-tracker .

RUN pwd
RUN ls -la

ENTRYPOINT /root/task-tracker
EXPOSE 8080