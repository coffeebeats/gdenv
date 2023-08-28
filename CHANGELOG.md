# Changelog

## 0.2.1 (2023-08-28)

## What's Changed
* feat(pkg/mirror): implement asset downloading by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/43
* chore(deps): bump tj-actions/changed-files from 37 to 38 by @dependabot in https://github.com/coffeebeats/gdenv/pull/45
* feat(cmd/gdenv): implement a `gdenv`-specific platform resolution function by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/46
* feat(pkg/godot): implement checksum operations `ExtractChecksum` and `ComputeChecksum` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/48


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.2.0...v0.2.1

## 0.2.0 (2023-08-26)

## What's Changed
* refactor!: migrate `cmd/gdenv`, `pkg/store`, and `pkg/pin` onto public `pkg/godot` package by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/41


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.1.4...v0.2.0

## 0.1.4 (2023-08-25)

## What's Changed
* feat(pkg/godot): implement a public `godot` package with a `Version` struct by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/38
* feat(pkg/godot): refactor `internal/godot` and add improved platform-handling logic by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/40


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.1.3...v0.1.4

## 0.1.3 (2023-08-12)

## What's Changed
* fix(ci): correctly identify release assets; use v-prefixed version tags in asset names by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/36


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.1.2...v0.1.3

## 0.1.2 (2023-08-12)

## What's Changed
* feat(gdenv): create skeleton implementations of `gdenv` and `gdenv-shim` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/27
* feat(gdenv): define flag options for all commands by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/29
* feat(internal/godot): implement a package with Godot specification functionality by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/31
* feat(pkg/pin): implement pin operations in `pkg/pin` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/32
* feat(pkg/store): implement core store functionality by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/33
* feat(gdenv/cmd): enable suggestions and short option handling by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/34
* feat(cmd/gdenv): implement more command functionality by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/35


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.1.1...v0.1.2

## 0.1.1 (2023-08-08)

## What's Changed
* chore(ci): remove pinned version in release workflow by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/18
* chore: add a `.gitattributes` file to handle line ending normalization by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/20
* chore: update issue templates for bugs and feature requests by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/21
* feat(docs): Add installation, usage, and meta sections to `README.md` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/22
* chore: add a PR template by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/23
* feat(ci): enable dependabot version updates by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/24
* feat(ci): add reviewers to dependabot PRs; check app deps daily by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/26
* chore(deps): bump golang.org/x/tools from 0.11.1 to 0.12.0 by @dependabot in https://github.com/coffeebeats/gdenv/pull/25

## New Contributors
* @dependabot made their first contribution in https://github.com/coffeebeats/gdenv/pull/25

**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.1.0...v0.1.1

## 0.1.0 (2023-08-08)

## What's Changed
* feat(ci): add a release workflow using `release-please` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/1
* feat(gdenv): create the `github.com/coffeebeats/gdenv` module by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/3
* feat(ci): set up a CI workflow `check-commit.yml` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/4
* feat(ci): set up application publishing using `goreleaser` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/5
* chore(ci): use github changelog type by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/12


**Full Changelog**: https://github.com/coffeebeats/gdenv/commits/v0.1.0
