# stage 1: build
FROM golang:1.16 as build
LABEL stage=intermediate
WORKDIR /app
COPY . .
RUN make build

# stage 2: scratch
FROM scratch as go-auth-service
EXPOSE 8080
COPY --from=build /app/bin/go-auth-service /bin/go-auth-service
USER 1000:1000
ENTRYPOINT ["/bin/go-auth-service"]