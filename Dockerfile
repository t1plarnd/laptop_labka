FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/laptop_labka .

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/laptop_labka /laptop_labka
EXPOSE 8080
ENTRYPOINT ["/laptop_labka"]
