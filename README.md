# Example Control Plane Function

Example function showcasing how you can use Crossplane Composition Functions to orchestrate your IDP's Control Plane

Needed Tools:
* [Crossplane CLI](https://docs.crossplane.io/latest/cli/#installing-the-cli)
* [IDPBuilder](https://cnoe.io/docs/reference-implementation/installations/idpbuilder)
* GoLang
* Kubernetes CLI
* Docker

To Test and Build Locally
```shell
# Run code generation - see input/generate.go
$ go generate ./...

# Run tests - see fn_test.go
$ go test ./...

# Build the function's runtime image - see Dockerfile
$ docker build . --tag=runtime

# Build a function package - see package/crossplane.yaml
$ crossplane xpkg build -f package --embed-runtime-image=runtime
```

To Test the output of the function based on a given Composite Resource
```shell
crossplane beta render ./example/xr.yaml ./example/composition.yaml ./example/functions.yaml
```

To Spin up your local IDP control plane pointing to a localstack instance.
```shell
idpbuilder create -p ./local/localstack
```

To apply a claim to validate the behavior of your Composition Function locally
```shell
kubectl apply -f example/claims/s3bucket.yaml
```