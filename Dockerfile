FROM golang:1.16 as builder
LABEL stage=intermediate
WORKDIR /app
COPY . .
RUN make build

FROM scratch as go-auth-service
EXPOSE 8080
COPY --from=builder /app/bin/go-auth-service /bin/go-auth-service
USER 1000:1000
ENTRYPOINT ["/bin/go-auth-service"]