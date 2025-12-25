# A simplistic URL shortener service

## Components

- cmd/service: Entrypoint of HTTP service.
- handlers: HTTP handlers based on `net/http.ServeMux`.
- storage: Storage URL and code, in JSON.

## Workflow

- Install Go 1.25+
- Start server with `go run ./cmd/service`
- Run tests with `go test -v -coverprofile=cov ./...`
- View test coverage with `go tool cover -html=coverage.out`
- Run linters with `golangci-lint run`
 
## Build and deploy to a local kind cluster

### Prerequisties
- [ ] Install [ko](https://github.com/google/ko)
- [ ] Install [kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [ ] Install [kubectl](https://kubernetes.io/docs/reference/kubectl/)
- [ ] Have a Kubernetes cluster running (e.g., kind)

### Use kind to build a cluster
MacOS can build a k8s cluster with below after [Homebrew](https://brew.sh/) is installed.
```shell
brew install ko kind kubernetes-cli
kind create cluster
```

### Use ko to build and push to kind cluster
Then build and push image to kind cluster
```shell
env KO_DOCKER_REPO=kind.local ko apply -f deploy/shortener.yml
kubectl get pods
kubectl wait --for=condition=Available deployments/shortener
```
After pods are ready, forward port to the service.
```shell
kubectl port-forward service/shortener 8080:8080
```
Then in other terminal window run `curl localhost:8080/ping` to test, use `kubectl logs -l app=shortener` to check logs.

After change some code, use the same command of build to rebuild and rollout pods.

Clean up with
```shell
ko delete -f deploy/shortener.yml
```

