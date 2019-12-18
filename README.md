# cloudkms-encryption

[![Go Report Card](https://goreportcard.com/badge/github.com/srinandan/cloundkms-encryption)](https://goreportcard.com/report/github.com/srinandan/cloudkms-encryption)

This service is meant to run as a sidecar to the Apigee Runtime (also known as Message Processor). The service takes a Google Cloud Service Account as a parameter and is used to encrypt or decrypt text using Cloud KMS.

## Use Case

This service is meant to be used with the Apigee hybrid [API Runtime](https://docs.apigee.com/hybrid). When developing API Proxies on Apigee, a developer may want to encrypt or decrypt parts of the payload. Google Cloud provides [Cloud KMS](https://cloud.google.com/kms/). The services uses Cloud KMS libraries to encrypt or decrpyt data.  

A typical pattern/example would be:

* Instantiate the `cloudkms-encryption` service as a service or as a sidecar to the Message Processor
* Use a [Service Callout policy](https://docs.apigee.com/api-platform/reference/policies/service-callout-policy) to first decide whether you want all or parts of the payload to be encrypted/decrypted. Then invoke the `cloudkms-encryption` service for encryption or decryption. 


## Usage

Input:

```bash

curl 0.0.0.0:8080/encrypt -d 'sample clear text data'
```

Output:

```bash

*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET /token HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Date: Wed, 18 Dec 2019 06:18:48 GMT
< Content-Length: 142
<
{"base64_cipher_text":"CiQATxZWh3Ky1nUed8+Uzfy1rrZ0hUrvt8J0OZUyauXbrvv2TwwSLwCPcW8BdQBpa9PXMWdOUk1c8SLNPG7J4NCyVXNfF8FLBnhgXYMGNCeY4B0673bf"}
```

### Install the sidecar

Modify the kubernetes [manifest](./cloudkms-encryption.yaml) and deploy it to Kubernetes. For example:

```bash

kubectl create secret -n {namespace} generic cloudkms-encryption-svc-account --from-file client_secret.json
kubectl apply -n {namespace} -f cloudkms-encryption.yaml
```
