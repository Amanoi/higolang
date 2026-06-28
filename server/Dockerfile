FROM golang:1.21-alpine AS build
RUN apk add --no-cache gcc musl-dev
WORKDIR /build
# Copy go.mod first, run tidy to generate go.sum inside container
COPY go.mod ./
RUN go mod tidy
# Copy source code and build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o higolang-server .

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /build/higolang-server .
COPY --from=build /build/config.yaml .
RUN mkdir -p data
EXPOSE 8080
CMD ["./higolang-server"]
