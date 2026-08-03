# longhand

Generated Go module (`github.com/Olian04/Longhand`, Go 1.26). Mode: `cli-library`.

## Layout

| Path | Role |
| --- | --- |
| `internal/domain/echo` | Demo domain model — the same in every mode. |
| `cmd/longhand` | CLI adapter: args/stdin → domain → stdout. |
| `pkg/longhand` | Public library facade over the domain. |
| `internal/config` | Config load (YAML → ENV → flags). |
| `configs/` | Example config files. |
| `internal/observability/logging` | slog setup + request-ID context helpers. |
| `test/unit/` | Unit tests mirroring `internal/` / `pkg/`. |

## Dependency direction

The domain is the fixed point; everything else adapts to it and depends inward.

`cmd` → domain, `pkg/longhand`, observability, config.
There is no `internal/app` in this mode; `cmd` is the composition root, and should
only parse input and wire dependencies — behaviour belongs in the domain.

## Config

Precedence: YAML → ENV → flags.
Point `APP_CONFIG_FILE` at a YAML file (omit for built-in defaults). See `configs/`.

Sections in this mode: `labels`, `logging`.

## Run locally

```bash
make help
make run
```

Reach the domain through the CLI. Logs go to stderr, so stdout stays pipeable:

```bash
make build
./dist/longhand hello world          # args
echo '  padded  ' | ./dist/longhand  # stdin
./dist/longhand hello 2>/dev/null    # result only
```

The library facade reaches the same domain the CLI does:

```go
import "github.com/Olian04/Longhand/pkg/longhand"

func main() {
	println(longhand.Echo("  hello  ")) // "hello"
}
```

## Checks

```bash
make format
make lint
make test
make test-race
```

## Releases
Push tag `v*`; CI runs GoReleaser for binaries.
Library release artifacts included for library modes.

## Agent orientation

See `docs/` and `CLAUDE.md` for agent context.
