# Release process

R2Lang is pre-1.0 and follows [Semantic Versioning](https://semver.org/):
`vMAJOR.MINOR.PATCH`.

- `MAJOR` stays `0` until the language/stdlib surface is considered stable.
- `MINOR` bumps for new features (new stdlib modules, new language syntax) —
  may still include breaking changes while `MAJOR` is `0`.
- `PATCH` bumps for bug fixes and hardening with no intentional behavior
  change to working code.

## Steps

1. Make sure `main` is green: `go build ./...`, `go vet ./...`,
   `go test ./... -count=1`, `go run main.go gold_test.r2`, and
   `./scripts/sweep_examples.sh` all pass.
2. Add an entry to the top of [CHANGELOG.md](CHANGELOG.md), under a new
   `## [MAJOR.MINOR.PATCH] - short title` heading, following the existing
   `Added`/`Changed`/`Fixed`/`Removed` structure. Add the corresponding
   compare-link reference at the bottom of the file.
3. Tag: `git tag -a vMAJOR.MINOR.PATCH -m "short summary"`.
4. Push: `git push origin main && git push origin vMAJOR.MINOR.PATCH`.

## Tag naming history

Tags up to and including `v0.1.25` used `vMAJOR.MINOR.PATCH_short-description`
(an informal shorthand used during a fast-moving overnight hardening sprint).
Starting from the next release, tags use plain `vMAJOR.MINOR.PATCH` — the
description lives in CHANGELOG.md and the annotated tag message instead.
Existing `_description`-suffixed tags are kept as-is; they are historical
record and are not renamed or deleted.
