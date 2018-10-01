# emon

Eventstore monitoring over HTTP

**NOTE: NOT ready for production!!!**

## Provider configuration

| Environment variable       | Description                     | Default               |
|----------------------------|---------------------------------|-----------------------|
| EMON_HTTP_BIND_ADDRESS     | The address emon will listen on | :8113                 |
| EMON_CLUSTER_HTTP_ENDPOINT | The eventstore HTTP address     | http://localhost:2113 |
| EMON_CLUSTER_SIZE          | The eventstore cluster size     | 3                     |




## Development

```bash
docker-compose up -d
./run
```

## Releasing

Pre-requisites:
- [goreleaser](https://goreleaser.com/)
- A github token

```bash
goreleaser --rm-dist
```
