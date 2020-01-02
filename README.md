# cloudkms-encryption

[![Go Report Card](https://goreportcard.com/badge/github.com/srinandan/cloudkms-encryption)](https://goreportcard.com/report/github.com/srinandan/cloudkms-encryption)

This service is meant to run as a sidecar to the Apigee hybrid API gateway (also known as Message Processor). The service takes a Google Cloud Service Account as a parameter and is used to encrypt or decrypt text using [Cloud KMS](https://cloud.google.com/kms/). The service can also store and retrieve data from GCP [Secret Manager](https://cloud.google.com/secret-manager/docs/).

## Use Case

This service is meant to be used with the Apigee hybrid [API Runtime](https://docs.apigee.com/hybrid). When developing API Proxies on Apigee, a developer may want to encrypt or decrypt parts of the payload. Google Cloud provides [Cloud KMS](https://cloud.google.com/kms/). The services uses Cloud KMS libraries to encrypt or decrpyt data.

Sensitive information often needs to be stored in a secure location. GCP Secret Manager provides a service (like a vault) to store sensitive information.

## Prerequisites

* Apigee hybrid runtime installed on GKE or GKE on-premises (v1.13.x)
* A GCP Project with Cloud KMS and Secret Manager APIs enabled
* A Service Account with the following roles:
  a. Cloud KMS CryptoKey Encrypter/Decrypter
  b. Secret Manager Admin
  c. Secret Manager Secret Accessor

## Prerequisites to build

* kubectl 1.13 or higher
* docker 19.x or higher
* skaffold 1.1.0 or higher (optional)

## Installation

### Installation via kubectl

1. Build the [docker image](./Dockerfile) `docker build -t gcr.io/{project-id}/cloudkms-encryption`
2. Push to a container registry `docker push gcr.io/{project-id}/cloudkms-encryption`
3. Modify the kubernetes [manifest](./cloudkms-encryption.yaml)

```bash

kubectl create secret -n {namespace} generic cloudkms-encryption-svc-account --from-file client_secret.json
kubectl apply -n {namespace} -f cloudkms-encryption.yaml
```

### Installation via Skaffold

This application can also be installed via [skaffold](https://skaffold.dev/). Modify the [skaffold.yaml](./skaffold.yaml) to set the appropriate project name.

```bash

skaffold run
```

#### Errors in Skaffold

When rerunning/installing the application, you may observe errors like this:

```bash

 - Error from server (Invalid): error when applying patch:
 ...
 ...
 `selector` does not match template `labels`
 ```

There is an open [issue](https://github.com/GoogleContainerTools/skaffold/issues/3133) for this in the skaffold project.

Workaound: first run `skaffold delete` and then `skaffold run`

## Supported Operations

### Encrypt data

Path: `/encrypt`
Method: `POST`
Accept: text/plain
Content-Type: application/json

The response is base64 encoded

```bash

curl 0.0.0.0:8080/encrypt -d 'sample clear text data'
```

Output:

```bash

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
<
{"payload":"CiQATxZWh3Ky1nUed8+Uzfy1rrZ0hUrvt8J0OZUyauXbrvv2TwwSLwCPcW8BdQBpa9PXMWdOUk1c8SLNPG7J4NCyVXNfF8FLBnhgXYMGNCeY4B0673bf"}
```

### Decrypt data

Path: `/decrypt`
Method: `POST`
Accept: text/plain
Content-Type: application/json 

```bash

curl 0.0.0.0:8080/decrypt -d 'CiQATxZWh3Ky1nUed8+Uzfy1rrZ0hUrvt8J0OZUyauXbrvv2TwwSLwCPcW8BdQBpa9PXMWdOUk1c8SLNPG7J4NCyVXNfF8FLBnhgXYMGNCeY4B0673bf'
```

Output:

```bash

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
<
{"payload":"sample clear text data"}
```

### Create Secret

Creates a new secret in Secret Manager.

Path: `/secrets`
Method: `POST`
Accept: application/json
Content-Type: application/json 

```bash

curl localhost:8080/secrets -H "Content-Type: application/json" -d '{"secretId":"test"}'
```

### Store Secret

Stores a secret in Secret Manager, optionally encrypts and stores a section in Secret Manager

Path: `/storesecrets`
Method: `POST`
Accept: application/json
Content-Type: application/json

```bash

curl localhost:8080/storesecrets -H "Content-Type: application/json" -d '{"secretId":"test","payload":"test data"}'
```

The same method can be used to encrypt first with Cloud KMS and then store in Secret Manager.

```json

{
  "secretId":"test",
  "payload":"test data",
  "encrypted": true
}
```

### Access a secret

Access a secret in Secret Manager, optionally decrypts the secret first and retrieves in clear text

Path: `/secrets/{secretName}/{version}`
Method: `GET`
Accept: application/json
Content-Type: application/json

```bash

curl localhost:8080/secrets/test/1
```

The same method can be used to access the data from Secret Manager and them decrypt with Cloud KMS

```bash

curl localhost:8080/secrets/test/1?ecrypted=true
```

## Access patterns from Apigee hyrid

A typical pattern/example would be to use a [Service Callout policy](https://docs.apigee.com/api-platform/reference/policies/service-callout-policy) to access operations supported by the service.


