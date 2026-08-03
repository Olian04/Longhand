# Configuration

Optional layered sources. **Precedence**: **YAML → ENV → flags** (later overrides earlier).

| Layer | Notes |
| ----- | ----- |
| YAML  | File via `Options.Path` / `APP_CONFIG_FILE` at the process entrypoint. Example: `configs/config.example.yaml`. |
| ENV   | `APP_LOGGING_*`. |
| Flags | `FlagOverrides` on `Options.Flags` (nil field = leave unchanged). |

Sections present in this mode: `labels`, `logging`.


With no layers enabled, `Load` returns validated defaults only.
