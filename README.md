faas-federation
-----

faas-federation is an implementation of the [faas-provider](https://github.com/openfaas/faas-provider) which can be used to unify one or more OpenFaaS clusters under a single API.

## Why do we need this?

This project exists to join together two or more distinct OpenFaaS clusters.

* Multi-region - east/west

Get a single API / control-plane for one or more clusters split by location, such as east/west.

* Edge locations

You may have one or more edge locations (or even ARM/IoT OpenFaaS installations). You can join them together under a single set of credentials and control plane.

* Varying provider types

You can connect two or more different OpenFaaS provider types together. For instance: Kubernetes (faas-netes) and Lambda (faas-lambda). This means you can have a single, centralized control-plane but deploy to both AWS Lambda and Kubernetes at the same time.

## Getting started

`faas-federation` can replace your provider in your existing OpenFaaS deployment.

More coming soon.

### Example

Coming soon: deploy OpenFaaS with two separate [`faas-memory`](https://github.com/openfaas-incubator/faas-memory) providers.

### helm chart

See also: example of Kubernetes and AWS Lambda federated configuration in the sample [helm chart](chart/of-federation).

## Gateway routing

To route to one gateway or another, simply set `com.openfaas.federation.gateway` to the name you want to pick.

| Annotation | Description |
| ----|----|
| `com.openfaas.federation.gateway` | route the request based on the provider name i.e. `faas-netes`, `faas-lambda` |

## Configuration

All configuration is managed using environment variables

| Option                            | Usage      | Default                  | Required |
|-----------------------------------|------------|--------------------------|----------|
| `providers`           | comma separated list of provider URLs i.e. `http://faas-netes:8080,http://faas-lambda:8080` | - |   yes    |
| `default_provider`    | default provider URLs used when no deployment constraints are matched i.e. `http://faas-netes:8080` | - |   yes    |

## Acknowledgements

Idea by Alex Ellis and Edward Wilde.

## License

MIT
