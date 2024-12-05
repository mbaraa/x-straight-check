FROM golang:1.22-alpine AS build

WORKDIR /app
COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest &&\
    apk add --no-cache make npm nodejs &&\
    make

FROM alpine:latest AS run

WORKDIR /app
COPY --from=build /app/x-straight-check .

EXPOSE 8080

CMD ["./x-straight-check"]
