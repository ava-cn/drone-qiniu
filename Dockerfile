# build go
FROM golang:1.15-alpine AS build
WORKDIR /go/src/qiniu
COPY . .
RUN go build -o /build/qiniu

FROM alpine
COPY --from=build /build/qiniu /bin/
ENTRYPOINT /bin/qiniu
