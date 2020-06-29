
#build stage
#FROM golang:alpine AS builder
FROM victron/exleacar_builder:latest AS builder
WORKDIR /go/src/exleacar
COPY . .
RUN apk add --no-cache git
RUN go get -d -v ./...
RUN go install -v ./...

#final stage
FROM alpine:latest
WORKDIR /home/ 
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/exleacar /usr/local/bin/
ENTRYPOINT ["exleacar", "-dir", "/home"]
CMD [ "-u", "", "-p", "",  "-w", "30", "-vvv"]
LABEL Name="exleacar" Version="2020-06-20"

# NOTES: 
# docker run --name exlea -d --rm --net=host -v /host_dir:/home --user $(id -u):$(id -g) victron/exleacar:2020-06-19_3 -u someOne -p secret -w 30 -vvv

