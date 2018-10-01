# emon

Eventstore monitoring over HTTP

**NOTE: NOT ready for production!!!**

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
