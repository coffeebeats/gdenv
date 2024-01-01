# Changelog

## 0.6.10 (2024-01-01)

## What's Changed
* feat(ci): add a workflow to auto-merge a Dependabot PR by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/183
* chore(deps): bump github.com/urfave/cli/v2 from 2.27.0 to 2.27.1 by @dependabot in https://github.com/coffeebeats/gdenv/pull/182
* fix(ci): remove example condition from workflow step by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/185


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.6.9...v0.6.10

## 0.6.9 (2023-12-31)

## What's Changed
* chore(deps): bump actions/setup-go from 4 to 5 by @dependabot in https://github.com/coffeebeats/gdenv/pull/171
* chore(deps): bump github/codeql-action from 2 to 3 by @dependabot in https://github.com/coffeebeats/gdenv/pull/173
* fix(ci): skip format job if triggered by dependabot by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/174
* chore: configure `markdownlint` to allow non-sibling repeat headings by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/175
* fix(ci): use correct dependabot name by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/180
* chore(deps): bump tj-actions/changed-files from 40 to 41 by @dependabot in https://github.com/coffeebeats/gdenv/pull/177
* chore(deps): bump github.com/urfave/cli/v2 from 2.26.0 to 2.27.0 by @dependabot in https://github.com/coffeebeats/gdenv/pull/178
* chore(deps): bump github.com/go-resty/resty/v2 from 2.10.0 to 2.11.0 by @dependabot in https://github.com/coffeebeats/gdenv/pull/179


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.6.8...v0.6.9

## 0.6.8 (2023-12-04)

## What's Changed
* chore(ci): migrate `release-please` to version `4` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/168
* fix(ci): correctly skip publish step if no release was created by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/170
* chore(deps): bump github.com/urfave/cli/v2 from 2.25.7 to 2.26.0 by @dependabot in https://github.com/coffeebeats/gdenv/pull/166


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.6.7...v0.6.8

## 0.6.7 (2023-12-03)

## What's Changed
* chore(deps): bump github.com/charmbracelet/log from 0.3.0 to 0.3.1 by @dependabot in https://github.com/coffeebeats/gdenv/pull/163
* chore: add exported environment variables to `.zshrc` instead of `.zshenv` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/165


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.6.6...v0.6.7

## 0.6.6 (2023-11-12)

## What's Changed
* chore(cmd/gdenv): remove extra newline after version by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/161


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.6.5...v0.6.6

## 0.6.5 (2023-11-12)

## What's Changed
* fix(cmd/gdenv-shim): ensure new process has Stdin connected by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/159


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.6.4...v0.6.5

## 0.6.4 (2023-11-10)

## What's Changed
* chore(docs): update link to compilation instructions by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/151
* fix(docs): correct Windows install script; simplify GitHub links in logging by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/156
* fix(tools): remove `tools.go` to simplify project dependencies by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/158


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.6.3...v0.6.4

## 0.6.3 (2023-11-08)

## What's Changed
* chore(docs): improve readability of `README.md` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/146
* chore: update `github.com/charmbracelet/log` to `v0.3.0` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/148
* chore: run `go mod tidy` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/150


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.6.2...v0.6.3

## 0.6.2 (2023-11-07)

## What's Changed
* chore(CI): use repository Go version during CodeQL scans by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/137
* chore(deps): bump actions/checkout from 3 to 4 by @dependabot in https://github.com/coffeebeats/gdenv/pull/132
* chore(deps): bump github.com/golangci/golangci-lint from 1.55.1 to 1.55.2 by @dependabot in https://github.com/coffeebeats/gdenv/pull/133
* chore(deps): bump github.com/goreleaser/goreleaser from 1.21.2 to 1.22.0 by @dependabot in https://github.com/coffeebeats/gdenv/pull/136
* chore(deps): bump golang.org/x/mod from 0.13.0 to 0.14.0 by @dependabot in https://github.com/coffeebeats/gdenv/pull/134
* fix(CI): run linting in separate job by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/139
* fix(CI): migrate to new `--skip` flag by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/140
* fix(CI): carry forward CodeCov coverage for entire project by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/141
* feat(CI): run tests with the race detector enabled by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/142
* feat(internal/fstest): create `Filepath` interface for creating different filepath types in tests by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/143
* fix(CI): ensure `go test` command has sufficient timeout by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/144
* fix(pkg/progress): eliminate deadlock in `TestWriter` test by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/145


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.6.1...v0.6.2

## 0.6.1 (2023-11-05)

## What's Changed
* chore(cmd/gdenv): remove `completions` command until improved support is added in `urfave/cli` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/130


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.6.0...v0.6.1

## 0.6.0 (2023-11-05)

## What's Changed
* chore: enabled `security-advanced` CodeQL queries by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/126
* refactor(pkg/godot/mirror)!: simplify `Mirror` usage by making it generic over `artifact.Versioned` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/128
* refactor(pkg/godot/artifact)!: simplify `artifact.Artifact` and `artifact/checksum.Checksums` implementations by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/129


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.5.3...v0.6.0

## 0.5.3 (2023-11-02)

## What's Changed
* feat(scripts): Add an `install.ps1` script for installing `gdenv` on Windows by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/122
* fix: update install instructions for PowerShell by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/124
* fix(pkg/godot/version): update `Version` to use `uint8` internally by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/125


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.5.2...v0.5.3

## 0.5.2 (2023-10-31)

## What's Changed
* feat(pkg/godot/version): add `GDENV_DEFAULT_MONO` to simplify Mono usage; improve `gdenv` version resolution logic by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/120
* chore(deps): bump github.com/docker/docker from 24.0.2+incompatible to 24.0.7+incompatible by @dependabot in https://github.com/coffeebeats/gdenv/pull/119


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.5.1...v0.5.2

## 0.5.1 (2023-10-30)

## What's Changed
* chore(docs): clean up `README.md` and split out `Commands` and `Installation` sections into `./docs` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/116
* chore(deps): bump google.golang.org/grpc from 1.57.0 to 1.57.1 by @dependabot in https://github.com/coffeebeats/gdenv/pull/118


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.5.0...v0.5.1

## 0.5.0 (2023-10-30)

## What's Changed
* feat: add support for tracking progress of installs and archive extraction by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/110
* refactor(internal/godot/mirror): split `Mirror` interface into separate concerns by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/112
* chore(deps): bump tj-actions/changed-files from 39 to 40 by @dependabot in https://github.com/coffeebeats/gdenv/pull/113
* feat(pkg/godot,pkg/progress)!: make `godot` and `progress` packages public by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/114
* fix(pkg/godot/version): ensure parsed integer has sufficient bit size by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/115


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.4.6...v0.5.0

## 0.4.6 (2023-10-28)

## What's Changed
* feat(cmd/gdenv): support installing source versions; add `vendor` command by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/107
* feat: add code coverage to pull requests via CodeCov by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/109


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.4.5...v0.4.6

## 0.4.5 (2023-10-25)

## What's Changed
* fix: reduce file permissions; require directory to exist for `pin.Write` by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/104
* fix(cmd/gdenv-shim): ensure first argument is binary name on mac/linux by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/106


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.4.4...v0.4.5

## 0.4.4 (2023-10-23)

## What's Changed
* fix: name the `gdenv-shim` binary `godot` to simplify installation by @coffeebeats in https://github.com/coffeebeats/gdenv/pull/102


**Full Changelog**: https://github.com/coffeebeats/gdenv/compare/v0.4.3...v0.4.4

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
