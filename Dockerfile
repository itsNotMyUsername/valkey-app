FROM golang:1.23.4-alpine3.21 as build

COPY ./main.go /app/main.go
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app
RUN go build -o /bin/app main.go

FROM scratch
COPY --from=build /bin/app /bin/app
CMD ["/bin/app"]
