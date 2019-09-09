FROM alpine:3.9

ADD . .

RUN apk update && apk add ca-certificates

EXPOSE 8110

CMD ["./main"]
