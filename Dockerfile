FROM golang:1.16.2-alpine3.13 AS build

WORKDIR /src

COPY . .

RUN go build -o main .

FROM 289208114389.dkr.ecr.us-east-1.amazonaws.com/alpine AS bin

COPY --from=build /src/main /

ENTRYPOINT ["/main"]
