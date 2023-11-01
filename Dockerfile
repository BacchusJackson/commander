FROM golang:alpine as build

WORKDIR /app
COPY . .

RUN apk --no-cache add ca-certificates git openssh make

RUN go build -ldflags="-w -s -extldflags '-static' -X 'main.Version=$(git rev-parse HEAD)'" -o /usr/local/bin/commander .

FROM alpine:latest as final-base

RUN apk --no-cache add curl jq sqlite-libs git ca-certificates tzdata

WORKDIR /app

LABEL org.opencontainers.image.title="Commander"
LABEL org.opencontainers.image.description="A simple CLI for templating Commands"
LABEL org.opencontainers.image.vendor="Bacchus Jackson"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL io.artifacthub.package.readme-url="https://raw.githubusercontent.com/BacchusJackson/commander/main/README.md"
LABEL io.artifacthub.package.license="Apache-2.0"

# Final image in a CI environment, assumes binaries are located in ./bin
# This is for pulling in prebuilt binaries and doesn't depend on the build job
FROM final-base as final-ci

COPY ./bin/commander /usr/local/bin/commander

# Final image if building locally and build dependencies are needed
FROM final-base

COPY --from=build /usr/local/bin/commander /usr/local/bin/commander

ENTRYPOINT [ "/usr/local/bin/commander" ]
