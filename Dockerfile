FROM alpine:latest

COPY app .
# this will copy the content of configs to root level at destination
COPY ./config.yaml .

EXPOSE 8080

CMD [ "./app" ]