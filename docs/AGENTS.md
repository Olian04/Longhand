# Agent context: longhand

App template — use with repo `README.md` and [`configs/config.example.yaml`](../configs/config.example.yaml). Mode: `cli-library`.

## Layout

| Path | Role |
| --- | --- |
| `internal/domain/echo` | Domain logic only. Identical in every mode; no IO imports. |


| `cmd/longhand` | CLI adapter over the domain: args/stdin in, stdout out. |


| `pkg/longhand` | Public API: exported facade delegating to the domain. |


| `internal/config` | Root config `Load`; section structs (defaults/`WithDefaults`/`Validate`). |
| `internal/observability/logging` | slog setup + request-ID context helpers. |
| `configs` | YAML contract examples. |


| `test/unit/...` | Unit tests beside mirrored paths. |

## Dependency direction

`internal/domain` is the fixed point. Mode-specific IO adapters depend on it; it
depends on nothing but the standard library. Put behaviour in the domain and
keep adapters to translation only.

`cmd` → `internal/domain/echo`, `pkg/longhand`, observability, aggregated `internal/config`.
There is no `internal/app` in this mode; `cmd` is the composition root.



## Configuration

`APP_CONFIG_FILE` optional; unset uses defaults matching example file shapes.

Section structs carry `Default*` constants in-package, `WithDefaults()`, `Validate()`. Root `config.Config.Validate()` delegates down.

Sections in this mode: `labels`, `logging`. Mode decides
which sections exist — do not add an `http` section to a mode with no server.


## Mode notes


- **Logging**: use `internal/observability/logging` for slog setup. Prefer
  `logging.FromContext(ctx)` over package-level `slog` so lines carry `request_id`.
- **Exit errors**: report from `main` to stderr, not `slog` — the configured
  logger is already torn down when the deferred cleanup has run.


- **Metrics**: off for this mode.


## Commands (`make`)

`lint`, `test`, `test-race`, `run`, `build` (output `./dist/longhand`).

---

## Go proverbs ([source](https://go-proverbs.github.io/), caveman compress)

Concurrency channels coordinate mutex serializes · not parallelism · small interface sharp · zero value useful · `any` untyped tame it · gofmt settles bikeshed · tiny copy beats dep hairball syscall/cgo build tags isolate · cgo not Go · `unsafe` no contract · clarity beats wit · reflection stay cold path · errors values inspect wrap once · architecture name docs users · panic stays in `main` / hard startup.

## Uber style distill ([guide](https://github.com/uber-go/guide/blob/master/style.md), caveman compress)

Rare `*Iface` · `var _ I = (*T)(nil)` at export boundary · defer unlock pairs · chan buffer zero or one usually · slice/map copy exported API boundaries · typed errors `%w` chain handle once · assert comma-ok · goroutine bounded ctx/waitgroup · no zombie `init()` · globals inject not mutate · exits from `main` only · strconv hot paths · structs field-named literals · table tests sub `t.Run`.

---

Canon links above beat bullet memory when tradeoff unclear.
