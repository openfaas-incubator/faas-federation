## Constraints
| Annotation | Description |
| ----|----|
| federation.provider_name | route the request based on the provider name i.e. `kubernetes`, `swarm` |

## Configuration
All configuration is managed using environment variables

| Option                            | Usage                                                                                          | Default                  | Required |
|-----------------------------------|------------------------------------------------------------------------------------------------|--------------------------|----------|
| `providers`           | comma separated list of provider URLs i.e. `http://faas-netes:8080,http://faas-lambda:8080` | - |   yes    |
| `default_provider`    | default provider URLs used when no deployment constraints are matched i.e. `http://faas-netes:8080` | - |   yes    |


