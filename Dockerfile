#build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git mercurial gcc
ADD . /src
RUN cd /src && go build -o service-discovery
LABEL stage=builder

#final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/service-discovery /app
RUN touch /app/targets.json /app/targets-ssl.json
ENTRYPOINT ./service-discovery