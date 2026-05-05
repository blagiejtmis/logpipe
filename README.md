# logpipe

Structured log aggregator that tails multiple sources and routes to configurable sinks.

## Installation

```bash
go install github.com/yourorg/logpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/logpipe.git && cd logpipe && go build ./...
```

## Usage

Define your sources and sinks in a config file:

```yaml
sources:
  - type: file
    path: /var/log/app/*.log
  - type: stdin

sinks:
  - type: elasticsearch
    url: http://localhost:9200
    index: logs
  - type: stdout
    format: json
```

Then run:

```bash
logpipe --config logpipe.yaml
```

You can also pipe logs directly:

```bash
tail -f /var/log/app.log | logpipe --sink stdout --format json
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Path to config file | `logpipe.yaml` |
| `--sink` | Output sink type | `stdout` |
| `--format` | Log output format | `json` |
| `--dry-run` | Parse config without running | `false` |

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

## License

[MIT](LICENSE)