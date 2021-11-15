![CI](https://github.com/jonascheng/stellar-upload-ratelimit/actions/workflows/ci.yaml/badge.svg)
![CD](https://github.com/jonascheng/stellar-upload-ratelimit/actions/workflows/cd.yaml/badge.svg)
![codecov](https://codecov.io/gh/jonascheng/stellar-upload-ratelimit/branch/main/graph/badge.svg)

# stellar-upload-ratelimit

## Usage

### Start simple upload server

```bash
$ docker-compose up
```

### Launch client to upload

```bash
$ ./bin/upload-go --help
usage: upload-go [<flags>]

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --server="http://localhost:8080/upload?token=secret"
             Upload server endpoint.
  --file="./data/agent-telemetry-threat-info-flat-500-1636954894.json.gz"
             Upload source file.
  --rate=5   Upload rate limit in Mbps.
  --version  Show application version.

```

## LICENSE

[MIT](https://github.com/jonascheng/stellar-upload-ratelimit/blob/master/LICENSE)
