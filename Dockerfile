# stage 1: build
FROM golang:1.15 as build
LABEL stage=intermediate
WORKDIR /app
COPY . .
RUN make build

# stage 2: scratch
FROM scratch as scratch
EXPOSE 8080
COPY --from=build /app/bin/go-auth-service /bin/go-auth-service
ENTRYPOINT ["go-auth-service"]