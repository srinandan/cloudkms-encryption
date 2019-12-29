#docker build -t gcr.io/srinandans-apigee/cloudkms-encryption .
#docker run -d -p 8080:8080 --name cloudkms-encryption -v ~/sa.json:/sa.json -e GOOGLE_APPLICATION_CREDENTIALS="./sa.json" -e DEBUG="true" gcr.io/nandanks-151422/cloudkms-encryption
FROM golang:latest as builder
RUN useradd -U app
ADD . /go/src/cloudkms-encryption 
WORKDIR /go/src/cloudkms-encryption
COPY . /go/src/cloudkms-encryption
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -extldflags "-static"' -o /go/bin/cloudkms-encryption

#without these certificates, we cannot verify the JWT token
FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
WORKDIR /
COPY --from=builder /go/bin/cloudkms-encryption .
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
USER app
CMD ["./cloudkms-encryption"]