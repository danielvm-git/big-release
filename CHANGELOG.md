# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

<!-- big-release managed -->

## [2.2.0] - 2026-07-24

### Added

- **e08:** semantic-release test parity ([#51](https://github.com/danielvm-git/big-release/issues/51)) ([736955e](https://github.com/danielvm-git/big-release/commit/736955ea14d6fa71f6d041a373c505740e98847d))

## [2.1.1] - 2026-07-24

### Fixed

- **ci:** remove redundant softprops upload step and enable verbose logger output ([9030e2e](https://github.com/danielvm-git/big-release/commit/9030e2e3e84086e5c84384e9449a079084170d92))

## [2.1.0] - 2026-07-24

### Added

- **pnpm:** add publisher with Detect Prepare Publish Verify ([350a0bb](https://github.com/danielvm-git/big-release/commit/350a0bbbfb079e52c2db789d8ce1e10dc2cf0ef7))

### Fixed

- **git:** use tagFormat in GetLastRelease ([4df1b30](https://github.com/danielvm-git/big-release/commit/4df1b306c298446d1a299eef883c14cd8210c52d))
- **config:** commit big-release config for CI ([a533279](https://github.com/danielvm-git/big-release/commit/a5332791780e5d31ca8052067a7bf4dd9d71d40a))
- **ci:** add debug to allowed commit types ([a110d44](https://github.com/danielvm-git/big-release/commit/a110d44334a8948848df9844f54901bebf1732d2))
- **ci:** allow merge commits in Conventional Commits check ([cad71cc](https://github.com/danielvm-git/big-release/commit/cad71cc5474d825d6518894572383543921b93c6))
- **ci:** configure git author for big-release commits ([1b1ef3f](https://github.com/danielvm-git/big-release/commit/1b1ef3f52c79681742ea5c65873d60becdad87e3))
- **ci:** use 'big-release release' subcommand ([22bfabc](https://github.com/danielvm-git/big-release/commit/22bfabc86c8a8930c23891ec3024c90ea1645d32))
- **ci:** remove tag gate from release workflow ([7a40792](https://github.com/danielvm-git/big-release/commit/7a407922edb61fce8bd3461fab53e184c81bffb7))

## [0.1.0] - 2026-07-24

### Added

- **pnpm:** add publisher with Detect Prepare Publish Verify ([350a0bb](https://github.com/danielvm-git/big-release/commit/350a0bbbfb079e52c2db789d8ce1e10dc2cf0ef7))
- **release:** channel management across git, npm, and GitHub (e22) ([#40](https://github.com/danielvm-git/big-release/issues/40)) ([3ccc86e](https://github.com/danielvm-git/big-release/commit/3ccc86ea3ceecebb5c9d89e039821503b4fabfb0))
- **plugins:** add GitLab release publisher plugin ([#39](https://github.com/danielvm-git/big-release/issues/39)) ([75a4da9](https://github.com/danielvm-git/big-release/commit/75a4da9ed2eba6370d54eb2e173ceb6b39c84473))
- **plugins:** add GitLab release publisher plugin (e23) ([c725b7d](https://github.com/danielvm-git/big-release/commit/c725b7d96eb41cdf03d53c4bbc5bea000cebe037))
- **plugins:** advanced GitHub and git release parity (e21) ([#38](https://github.com/danielvm-git/big-release/issues/38)) ([3ceef96](https://github.com/danielvm-git/big-release/commit/3ceef96abe883dd9f65cf53a9cd9d0b619a5fc69))
- **plugins:** advanced GitHub and git release parity (e21) ([#36](https://github.com/danielvm-git/big-release/issues/36)) ([a46f1a0](https://github.com/danielvm-git/big-release/commit/a46f1a0ef13e3e6dff9dca2606dff9249e630857))
- **release:** multi-branch support with glob matching and channels (e20,[#14](https://github.com/danielvm-git/big-release/issues/14)) ([46554a1](https://github.com/danielvm-git/big-release/commit/46554a10c8bdda8f1b5e3aa5e5d007b0527bacc5))
- **github:** comment on resolved issues/PRs after release (e19s04,[#12](https://github.com/danielvm-git/big-release/issues/12)) ([e42793c](https://github.com/danielvm-git/big-release/commit/e42793ca9158e3f4852f42cc9b3d0f834e9cf0e9))
- **github:** configurable release name and body templates (e19s03,[#11](https://github.com/danielvm-git/big-release/issues/11)) ([9c8b39f](https://github.com/danielvm-git/big-release/commit/9c8b39f3dfef05da1938b9030423d853e679391a))
- **github:** support draft GitHub releases (e19s02,[#13](https://github.com/danielvm-git/big-release/issues/13)) ([2370522](https://github.com/danielvm-git/big-release/commit/2370522d3757181c37914578ec4919f513cb24a3))
- **notes:** filter revert commits from release notes (e18s04,[#9](https://github.com/danielvm-git/big-release/issues/9)) ([49730e4](https://github.com/danielvm-git/big-release/commit/49730e4f9ba5fb7af7179c276e5bba67c110e08d))
- **notes:** clickable commit hash and issue links (e18s03,[#8](https://github.com/danielvm-git/big-release/issues/8)) ([b0d20cf](https://github.com/danielvm-git/big-release/commit/b0d20cf6be8599a3df3a39255ebb0762209328ce))
- **changelog:** hide non-release commit types by default (e18s02,[#7](https://github.com/danielvm-git/big-release/issues/7)) ([193cec2](https://github.com/danielvm-git/big-release/commit/193cec2d6a4660467e8a83e4f365821f0aad3817))
- **algorithm:** configurable commit type sections and visibility (e18s01,[#15](https://github.com/danielvm-git/big-release/issues/15)) ([da5ab53](https://github.com/danielvm-git/big-release/commit/da5ab536ce4d7b1e2d9e6ec6dc61da055c6ba5a1))
- **feat:** split Context into ReadOnlyContext and MutableState ([#5](https://github.com/danielvm-git/big-release/issues/5)) ([5e0fd0e](https://github.com/danielvm-git/big-release/commit/5e0fd0e73cae83e56e13b5dcd0228a3e0e12f523))
- **e06:** Algorithm & Plugin Unification ([#4](https://github.com/danielvm-git/big-release/issues/4)) ([77c4c21](https://github.com/danielvm-git/big-release/commit/77c4c2130cc580799307a300be042fa34340fc5d))
- **publishers,git:** e05 Publisher & Git Hardening ([#3](https://github.com/danielvm-git/big-release/issues/3)) ([a6f4ccb](https://github.com/danielvm-git/big-release/commit/a6f4ccbb3228f9bf4aae6e779e404ec2f6babb38))
- **release:** add semantic release parity — initial version, version calc, branch validation, CI detection, VerifyRelease step ([1a487b6](https://github.com/danielvm-git/big-release/commit/1a487b6ccf546bcc07b1384200b8b81be5a6c52a))
- **e03:** implement plugin system with git, github, exec, and changelog plugins ([2114744](https://github.com/danielvm-git/big-release/commit/2114744940add60463a0b7d14e0bbf494906bae3))
- **e02:** add Godot publisher ([3e2c4a1](https://github.com/danielvm-git/big-release/commit/3e2c4a12d137c914e9ed6e2626ba59ebccd7c3fd))
- **e02:** add Swift publisher audit ([ce3537e](https://github.com/danielvm-git/big-release/commit/ce3537ee7c10b38f11b4529a26ebf616c6d0b27d))
- **e02:** add Maven publisher audit ([abaf7cb](https://github.com/danielvm-git/big-release/commit/abaf7cbc6ed71461f9e462d9bdeda05e22002c89))
- **e02:** add Go Proxy publisher audit ([c191893](https://github.com/danielvm-git/big-release/commit/c19189353cf683a3588d19a672542cc473ab7648))
- **e02:** add crates.io publisher audit ([d39043f](https://github.com/danielvm-git/big-release/commit/d39043f94d0003ad50c1c565377d81f04b47e1ab))
- **e02:** add PyPI publisher audit ([d3c86fe](https://github.com/danielvm-git/big-release/commit/d3c86fe99a86162c35ae127f32a2a3afaa217633))
- **publishers:** add Maven publisher ([fe612ac](https://github.com/danielvm-git/big-release/commit/fe612ac4671267d3559e980db60bae17481d18e6))
- **publishers:** add Packagist publisher ([459fe3d](https://github.com/danielvm-git/big-release/commit/459fe3dd3205c9899bf81816ad0e11e41ad30014))
- **publishers:** add Swift publisher ([c4c7d6a](https://github.com/danielvm-git/big-release/commit/c4c7d6a6de37f010018582091c712c76d57a29b4))
- **publishers:** add Go Proxy publisher ([ca7bf55](https://github.com/danielvm-git/big-release/commit/ca7bf558e875d4ccddef37c0f2b6adab015667c4))
- **publishers:** add crates.io publisher ([82bf4ec](https://github.com/danielvm-git/big-release/commit/82bf4eca0e19c9537d5da6a0bcf1bd9193a03522))
- **publishers:** add PyPI publisher ([15ded65](https://github.com/danielvm-git/big-release/commit/15ded652eb05c997a22194feee3752570b250053))
- **feat:** add health command, guard-git hooks, CI workflows, and fix build issues ([dc25e19](https://github.com/danielvm-git/big-release/commit/dc25e1958afc6fb3eb0c6c646fdc8d2b584fb3c2))
- **feat:** initialize big-release project structure ([ef778d3](https://github.com/danielvm-git/big-release/commit/ef778d3e6cfa0ffbf6c1cde98f15763ae71e601e))

### Fixed

- **config:** commit big-release config for CI ([a533279](https://github.com/danielvm-git/big-release/commit/a5332791780e5d31ca8052067a7bf4dd9d71d40a))
- **ci:** add debug to allowed commit types ([a110d44](https://github.com/danielvm-git/big-release/commit/a110d44334a8948848df9844f54901bebf1732d2))
- **ci:** allow merge commits in Conventional Commits check ([cad71cc](https://github.com/danielvm-git/big-release/commit/cad71cc5474d825d6518894572383543921b93c6))
- **ci:** configure git author for big-release commits ([1b1ef3f](https://github.com/danielvm-git/big-release/commit/1b1ef3f52c79681742ea5c65873d60becdad87e3))
- **ci:** use 'big-release release' subcommand ([22bfabc](https://github.com/danielvm-git/big-release/commit/22bfabc86c8a8930c23891ec3024c90ea1645d32))
- **ci:** remove tag gate from release workflow ([7a40792](https://github.com/danielvm-git/big-release/commit/7a407922edb61fce8bd3461fab53e184c81bffb7))
- **ci:** run attribution check on pull requests only ([#43](https://github.com/danielvm-git/big-release/issues/43)) ([3b104cc](https://github.com/danielvm-git/big-release/commit/3b104cc1289e27f6a42c142f11c5fe74a6e7ad1e))
- **ci,config:** unblock release verify and default config validation ([#42](https://github.com/danielvm-git/big-release/issues/42)) ([dd090e8](https://github.com/danielvm-git/big-release/commit/dd090e846ee2a74fe028803aef35f32e95539d96))
- **release:** skip release on pull request CI builds ([#33](https://github.com/danielvm-git/big-release/issues/33)) ([aa82522](https://github.com/danielvm-git/big-release/commit/aa82522a3df49a19206f23c2ccfa7ee5b65b7832))
- **changelog:** Keep-a-Changelog 1.1.0 format and configurable title ([#32](https://github.com/danielvm-git/big-release/issues/32)) ([a45b636](https://github.com/danielvm-git/big-release/commit/a45b63697390af3333d880bc3b5820bab0fb23ea))
- **ci:** ship platform binaries and bump checkout to v7 ([#30](https://github.com/danielvm-git/big-release/issues/30)) ([f252013](https://github.com/danielvm-git/big-release/commit/f25201309d0e86e98a0e2ac8d478638c7e57847e))
- **config:** reject duplicate branch names in ValidateConfig ([#29](https://github.com/danielvm-git/big-release/issues/29)) ([a8b7f8b](https://github.com/danielvm-git/big-release/commit/a8b7f8b328bd4a71fd6f79a6888c32c702a16ed6))
- **config:** clean up conflict marker ([b520356](https://github.com/danielvm-git/big-release/commit/b5203564629fae17f64723f41e580e8d84b8e218))
- **release:** wire Analyzer fallback and nil guard in resolveNotes ([#6](https://github.com/danielvm-git/big-release/issues/6)) ([1c2453d](https://github.com/danielvm-git/big-release/commit/1c2453dbc7b383f9d7bd82aed931d3d06552e3bd))
- **release:** propagate branch config and respect publishers.enabled flag ([#2](https://github.com/danielvm-git/big-release/issues/2)) ([1d7bb19](https://github.com/danielvm-git/big-release/commit/1d7bb1903f3773f3601d53a5864b0018762a9326))
- **cli:** harden CLI error handling, publisher dry-run, and orchestrator package ([0491dcb](https://github.com/danielvm-git/big-release/commit/0491dcb692b05137505f475ae50af827f9bbc2f3))
- **ci:** use --format=%s instead of --oneline in CC check ([caea50b](https://github.com/danielvm-git/big-release/commit/caea50b15b1ee2fe112729116b6127d02338bc16))
- **ci:** bump golangci-lint-action to v7 for Go 1.26 support ([3ab116f](https://github.com/danielvm-git/big-release/commit/3ab116fdb655e5acc35288bb94171638e7131ed8))
- **e03:** respond to review — tab splitting and orphan tag cleanup ([1845229](https://github.com/danielvm-git/big-release/commit/184522966337710d4495a3eed1ffccb66fbff575))
- **e03:** address review findings — exec quoting, changelog safety, git test ([6035983](https://github.com/danielvm-git/big-release/commit/6035983c31a3fa3c4bb1327dcf956b539a03602c))
- **ci:** resolve release workflow failures and bump Go version ([073b135](https://github.com/danielvm-git/big-release/commit/073b1358eac4db2d08c5c8f60aee2ebafa4f99b3))
- **ci:** install golangci-lint in verify job of release.yml ([d2421b2](https://github.com/danielvm-git/big-release/commit/d2421b207dd987fe33b7b3de8099478cea4d8a96))
- **ci:** add missing pkg/release package and fix .gitignore ([04d0793](https://github.com/danielvm-git/big-release/commit/04d07934a3df7bfb62ea04f3cab8193f3618075f))
- **e02:** add npm publisher tests and address review findings ([838ce86](https://github.com/danielvm-git/big-release/commit/838ce86d5c998804d022ed89bb9f9bc50d5970fb))

## [0.1.0] - 2026-07-24

### Added

- **pnpm:** add publisher with Detect Prepare Publish Verify ([350a0bb](https://github.com/danielvm-git/big-release/commit/350a0bbbfb079e52c2db789d8ce1e10dc2cf0ef7))
- **release:** channel management across git, npm, and GitHub (e22) ([#40](https://github.com/danielvm-git/big-release/issues/40)) ([3ccc86e](https://github.com/danielvm-git/big-release/commit/3ccc86ea3ceecebb5c9d89e039821503b4fabfb0))
- **plugins:** add GitLab release publisher plugin ([#39](https://github.com/danielvm-git/big-release/issues/39)) ([75a4da9](https://github.com/danielvm-git/big-release/commit/75a4da9ed2eba6370d54eb2e173ceb6b39c84473))
- **plugins:** add GitLab release publisher plugin (e23) ([c725b7d](https://github.com/danielvm-git/big-release/commit/c725b7d96eb41cdf03d53c4bbc5bea000cebe037))
- **plugins:** advanced GitHub and git release parity (e21) ([#38](https://github.com/danielvm-git/big-release/issues/38)) ([3ceef96](https://github.com/danielvm-git/big-release/commit/3ceef96abe883dd9f65cf53a9cd9d0b619a5fc69))
- **plugins:** advanced GitHub and git release parity (e21) ([#36](https://github.com/danielvm-git/big-release/issues/36)) ([a46f1a0](https://github.com/danielvm-git/big-release/commit/a46f1a0ef13e3e6dff9dca2606dff9249e630857))
- **release:** multi-branch support with glob matching and channels (e20,[#14](https://github.com/danielvm-git/big-release/issues/14)) ([46554a1](https://github.com/danielvm-git/big-release/commit/46554a10c8bdda8f1b5e3aa5e5d007b0527bacc5))
- **github:** comment on resolved issues/PRs after release (e19s04,[#12](https://github.com/danielvm-git/big-release/issues/12)) ([e42793c](https://github.com/danielvm-git/big-release/commit/e42793ca9158e3f4852f42cc9b3d0f834e9cf0e9))
- **github:** configurable release name and body templates (e19s03,[#11](https://github.com/danielvm-git/big-release/issues/11)) ([9c8b39f](https://github.com/danielvm-git/big-release/commit/9c8b39f3dfef05da1938b9030423d853e679391a))
- **github:** support draft GitHub releases (e19s02,[#13](https://github.com/danielvm-git/big-release/issues/13)) ([2370522](https://github.com/danielvm-git/big-release/commit/2370522d3757181c37914578ec4919f513cb24a3))
- **notes:** filter revert commits from release notes (e18s04,[#9](https://github.com/danielvm-git/big-release/issues/9)) ([49730e4](https://github.com/danielvm-git/big-release/commit/49730e4f9ba5fb7af7179c276e5bba67c110e08d))
- **notes:** clickable commit hash and issue links (e18s03,[#8](https://github.com/danielvm-git/big-release/issues/8)) ([b0d20cf](https://github.com/danielvm-git/big-release/commit/b0d20cf6be8599a3df3a39255ebb0762209328ce))
- **changelog:** hide non-release commit types by default (e18s02,[#7](https://github.com/danielvm-git/big-release/issues/7)) ([193cec2](https://github.com/danielvm-git/big-release/commit/193cec2d6a4660467e8a83e4f365821f0aad3817))
- **algorithm:** configurable commit type sections and visibility (e18s01,[#15](https://github.com/danielvm-git/big-release/issues/15)) ([da5ab53](https://github.com/danielvm-git/big-release/commit/da5ab536ce4d7b1e2d9e6ec6dc61da055c6ba5a1))
- **feat:** split Context into ReadOnlyContext and MutableState ([#5](https://github.com/danielvm-git/big-release/issues/5)) ([5e0fd0e](https://github.com/danielvm-git/big-release/commit/5e0fd0e73cae83e56e13b5dcd0228a3e0e12f523))
- **e06:** Algorithm & Plugin Unification ([#4](https://github.com/danielvm-git/big-release/issues/4)) ([77c4c21](https://github.com/danielvm-git/big-release/commit/77c4c2130cc580799307a300be042fa34340fc5d))
- **publishers,git:** e05 Publisher & Git Hardening ([#3](https://github.com/danielvm-git/big-release/issues/3)) ([a6f4ccb](https://github.com/danielvm-git/big-release/commit/a6f4ccbb3228f9bf4aae6e779e404ec2f6babb38))
- **release:** add semantic release parity — initial version, version calc, branch validation, CI detection, VerifyRelease step ([1a487b6](https://github.com/danielvm-git/big-release/commit/1a487b6ccf546bcc07b1384200b8b81be5a6c52a))
- **e03:** implement plugin system with git, github, exec, and changelog plugins ([2114744](https://github.com/danielvm-git/big-release/commit/2114744940add60463a0b7d14e0bbf494906bae3))
- **e02:** add Godot publisher ([3e2c4a1](https://github.com/danielvm-git/big-release/commit/3e2c4a12d137c914e9ed6e2626ba59ebccd7c3fd))
- **e02:** add Swift publisher audit ([ce3537e](https://github.com/danielvm-git/big-release/commit/ce3537ee7c10b38f11b4529a26ebf616c6d0b27d))
- **e02:** add Maven publisher audit ([abaf7cb](https://github.com/danielvm-git/big-release/commit/abaf7cbc6ed71461f9e462d9bdeda05e22002c89))
- **e02:** add Go Proxy publisher audit ([c191893](https://github.com/danielvm-git/big-release/commit/c19189353cf683a3588d19a672542cc473ab7648))
- **e02:** add crates.io publisher audit ([d39043f](https://github.com/danielvm-git/big-release/commit/d39043f94d0003ad50c1c565377d81f04b47e1ab))
- **e02:** add PyPI publisher audit ([d3c86fe](https://github.com/danielvm-git/big-release/commit/d3c86fe99a86162c35ae127f32a2a3afaa217633))
- **publishers:** add Maven publisher ([fe612ac](https://github.com/danielvm-git/big-release/commit/fe612ac4671267d3559e980db60bae17481d18e6))
- **publishers:** add Packagist publisher ([459fe3d](https://github.com/danielvm-git/big-release/commit/459fe3dd3205c9899bf81816ad0e11e41ad30014))
- **publishers:** add Swift publisher ([c4c7d6a](https://github.com/danielvm-git/big-release/commit/c4c7d6a6de37f010018582091c712c76d57a29b4))
- **publishers:** add Go Proxy publisher ([ca7bf55](https://github.com/danielvm-git/big-release/commit/ca7bf558e875d4ccddef37c0f2b6adab015667c4))
- **publishers:** add crates.io publisher ([82bf4ec](https://github.com/danielvm-git/big-release/commit/82bf4eca0e19c9537d5da6a0bcf1bd9193a03522))
- **publishers:** add PyPI publisher ([15ded65](https://github.com/danielvm-git/big-release/commit/15ded652eb05c997a22194feee3752570b250053))
- **feat:** add health command, guard-git hooks, CI workflows, and fix build issues ([dc25e19](https://github.com/danielvm-git/big-release/commit/dc25e1958afc6fb3eb0c6c646fdc8d2b584fb3c2))
- **feat:** initialize big-release project structure ([ef778d3](https://github.com/danielvm-git/big-release/commit/ef778d3e6cfa0ffbf6c1cde98f15763ae71e601e))

### Fixed

- **ci:** configure git author for big-release commits ([1b1ef3f](https://github.com/danielvm-git/big-release/commit/1b1ef3f52c79681742ea5c65873d60becdad87e3))
- **ci:** use 'big-release release' subcommand ([22bfabc](https://github.com/danielvm-git/big-release/commit/22bfabc86c8a8930c23891ec3024c90ea1645d32))
- **ci:** remove tag gate from release workflow ([7a40792](https://github.com/danielvm-git/big-release/commit/7a407922edb61fce8bd3461fab53e184c81bffb7))
- **ci:** run attribution check on pull requests only ([#43](https://github.com/danielvm-git/big-release/issues/43)) ([3b104cc](https://github.com/danielvm-git/big-release/commit/3b104cc1289e27f6a42c142f11c5fe74a6e7ad1e))
- **ci,config:** unblock release verify and default config validation ([#42](https://github.com/danielvm-git/big-release/issues/42)) ([dd090e8](https://github.com/danielvm-git/big-release/commit/dd090e846ee2a74fe028803aef35f32e95539d96))
- **release:** skip release on pull request CI builds ([#33](https://github.com/danielvm-git/big-release/issues/33)) ([aa82522](https://github.com/danielvm-git/big-release/commit/aa82522a3df49a19206f23c2ccfa7ee5b65b7832))
- **changelog:** Keep-a-Changelog 1.1.0 format and configurable title ([#32](https://github.com/danielvm-git/big-release/issues/32)) ([a45b636](https://github.com/danielvm-git/big-release/commit/a45b63697390af3333d880bc3b5820bab0fb23ea))
- **ci:** ship platform binaries and bump checkout to v7 ([#30](https://github.com/danielvm-git/big-release/issues/30)) ([f252013](https://github.com/danielvm-git/big-release/commit/f25201309d0e86e98a0e2ac8d478638c7e57847e))
- **config:** reject duplicate branch names in ValidateConfig ([#29](https://github.com/danielvm-git/big-release/issues/29)) ([a8b7f8b](https://github.com/danielvm-git/big-release/commit/a8b7f8b328bd4a71fd6f79a6888c32c702a16ed6))
- **config:** clean up conflict marker ([b520356](https://github.com/danielvm-git/big-release/commit/b5203564629fae17f64723f41e580e8d84b8e218))
- **release:** wire Analyzer fallback and nil guard in resolveNotes ([#6](https://github.com/danielvm-git/big-release/issues/6)) ([1c2453d](https://github.com/danielvm-git/big-release/commit/1c2453dbc7b383f9d7bd82aed931d3d06552e3bd))
- **release:** propagate branch config and respect publishers.enabled flag ([#2](https://github.com/danielvm-git/big-release/issues/2)) ([1d7bb19](https://github.com/danielvm-git/big-release/commit/1d7bb1903f3773f3601d53a5864b0018762a9326))
- **cli:** harden CLI error handling, publisher dry-run, and orchestrator package ([0491dcb](https://github.com/danielvm-git/big-release/commit/0491dcb692b05137505f475ae50af827f9bbc2f3))
- **ci:** use --format=%s instead of --oneline in CC check ([caea50b](https://github.com/danielvm-git/big-release/commit/caea50b15b1ee2fe112729116b6127d02338bc16))
- **ci:** bump golangci-lint-action to v7 for Go 1.26 support ([3ab116f](https://github.com/danielvm-git/big-release/commit/3ab116fdb655e5acc35288bb94171638e7131ed8))
- **e03:** respond to review — tab splitting and orphan tag cleanup ([1845229](https://github.com/danielvm-git/big-release/commit/184522966337710d4495a3eed1ffccb66fbff575))
- **e03:** address review findings — exec quoting, changelog safety, git test ([6035983](https://github.com/danielvm-git/big-release/commit/6035983c31a3fa3c4bb1327dcf956b539a03602c))
- **ci:** resolve release workflow failures and bump Go version ([073b135](https://github.com/danielvm-git/big-release/commit/073b1358eac4db2d08c5c8f60aee2ebafa4f99b3))
- **ci:** install golangci-lint in verify job of release.yml ([d2421b2](https://github.com/danielvm-git/big-release/commit/d2421b207dd987fe33b7b3de8099478cea4d8a96))
- **ci:** add missing pkg/release package and fix .gitignore ([04d0793](https://github.com/danielvm-git/big-release/commit/04d07934a3df7bfb62ea04f3cab8193f3618075f))
- **e02:** add npm publisher tests and address review findings ([838ce86](https://github.com/danielvm-git/big-release/commit/838ce86d5c998804d022ed89bb9f9bc50d5970fb))
