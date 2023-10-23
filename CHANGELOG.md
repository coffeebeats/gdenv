# Changelog

## 0.4.3 (2023-10-23)

## What's Changed
* chore(deps): bump github.com/charmbracelet/lipgloss from 0.8.0 to 0.9.1 by @dependabot in https://github.com/coffeebeats/gdenv/pull/99
* fix(cmd/gdenv-shim,scripts): add Windows `bash` support to install script and `gdenv-shim` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/101


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.4.2...v0.4.3

## 0.4.2 (2023-10-23)

## What's Changed
* feat: add logging throughout the application by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/96
* feat(scripts): add an `sh`-compatible install script by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/98


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.4.1...v0.4.2

## 0.4.1 (2023-10-22)

## What's Changed
* fix(internal/godot/mirror): return error if no `Mirror` found by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/93
* refactor(pkg/store,pkg/pin): simplify API and improve test coverage by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/95


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.4.0...v0.4.1

## 0.4.0 (2023-10-16)

## What's Changed
* feat(pkg/artifact): implement a `Folder` artifact and update consumers by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/73
* refactor(internal/godot/artifacts)!: simplify `artifacts` package; move `mirror` under `godot` package by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/76
* chore(deps): bump github.com/go-resty/resty/v2 from 2.8.0 to 2.9.1 by @dependabot in https://github.com/coffeebeats/gdenv/pull/75
* feat(internal/godot/artifact): implement source and executable archive extraction by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/80
* chore(deps): bump golang.org/x/mod from 0.12.0 to 0.13.0 by @dependabot in https://github.com/coffeebeats/gdenv/pull/77
* chore(deps): bump golang.org/x/tools from 0.13.0 to 0.14.0 by @dependabot in https://github.com/coffeebeats/gdenv/pull/78
* chore(deps): bump golang.org/x/net from 0.15.0 to 0.17.0 by @dependabot in https://github.com/coffeebeats/gdenv/pull/79
* refactor(internal/godot/mirror): simplify `mirror` method `ExecutableArchive` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/81
* fix: use correct `fs.FileMode` when writing files and directories by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/82
* feat(internal/godot/mirror): implement a `mirror.Choose` function; utilize `context.Context` in `internal/client` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/83
* feat: propagate `context.Context` throughout application; improve CLI exit handling by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/84
* chore: increase cyclomatic complexity limit to `12` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/85
* feat(internal/godot/platform): define a `platform.Detect` function for resolving the target install platform by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/86
* feat(pkg/store): define new `ExecutePath`; correct `ToolPath` implementation by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/87
* fix(internal/godot/artifact/executable): ensure macOS executable path includes OS-appropriate separators by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/88
* feat(pkg/download): implement functions to download artifacts by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/89
* feat(pkg/install): implement full installation functionality for source and executables by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/90
* feat(cmd/gdenv-shim): implement the shim executable by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/92
* chore(deps): bump github.com/go-resty/resty/v2 from 2.9.1 to 2.10.0 by @dependabot in https://github.com/coffeebeats/gdenv/pull/91


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.3.3...v0.4.0

## 0.3.3 (2023-09-18)

## What's Changed
* feat(internal/godot/artifact,pkg/godot): create new `artifact` package; remove `pkg/godot` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/69
* chore(deps): bump github.com/go-resty/resty/v2 from 2.7.0 to 2.8.0 by @dependabot in https://github.com/coffeebeats/gdenv/pull/68
* chore(deps): bump goreleaser/goreleaser-action from 4 to 5 by @dependabot in https://github.com/coffeebeats/gdenv/pull/67


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.3.2...v0.3.3

## 0.3.2 (2023-09-16)

## What's Changed
* refactor(pkg/godot,internal/version): move the `Version` implementation to separate package by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/63
* chore(deps): bump actions/checkout from 3 to 4 by @dependabot in https://github.com/coffeebeats/gdenv/pull/61
* chore(deps): bump tj-actions/changed-files from 38 to 39 by @dependabot in https://github.com/coffeebeats/gdenv/pull/62
* refactor(pkg/godot,internal/platform): separate `Platform` into internal package `platform` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/65
* chore(internal/godot): move `platform` and `version` packages under `internal/godot` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/66


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.3.1...v0.3.2

## 0.3.1 (2023-09-09)

## What's Changed
* refactor: numerous minor refactors to be more idiomatic/improve readability by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/55
* feat(internal/progress): update `Progress` API to enable post-initialization configuration by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/57
* chore(deps): bump golang.org/x/tools from 0.12.0 to 0.13.0 by @dependabot in https://github.com/coffeebeats/gdenv/pull/58
* feat(internal/client,pkg/mirror): add a `Client.Exists` method; add a `Mirror.Has` method by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/59
* fix(pkg/godot): make `Platform` usage safer by restricting visibility for fields by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/60


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.3.0...v0.3.1

## 0.3.0 (2023-09-03)

## What's Changed
* feat: add `main` as a protected branch in vs code by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/49
* fix(pkg/godot): improve platform handling, especially for `mono` builds by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/51
* feat(pkg/mirror): improve the `mirror` package by factoring out client logic by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/52
* feat(pkg/progress): create `Progress` and `progress.Writer` structs for tracking progress by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/53


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.2.1...v0.3.0

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
