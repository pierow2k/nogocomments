<!-- markdownlint-disable no-duplicate-heading -->
# nogocomments Changelog

All notable changes to nogocomments will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Removed

## [3.0.0] - 2026-03-24

### Breaking

- Removed the --debug / -d CLI flag.
- Removed the pkg/filereader package from the repository.

### Changed

- Refactored the root command to inline file and clipboard input handling using os.ReadFile and clipboard.ReadAll().
- Removed the logrus dependency from the CLI path.
- Slightly clarified the --paste flag help text.
- Improved inline code documentation and comments throughout cmd/root.go.
- Update GitHub actions.

## [2.0.0] - 2026-03-20

- Bump Go to v1.26.1
- Refactor to use cobra for CLI
- Remove legacy internal packages with updated packages
- Update documentation

## [1.0.3] - 2025-07-25

- Version 1.0.3 uses Go version 1.24.5 but is otherwise functionally identical to prior versions.

## [1.0.0] - 2026-03-22

### Added

- Initial release.

[unreleased]: https://github.com/pierow2k/nogocomments/compare/v1.1.0...HEAD
[3.0.0]: https://github.com/pierow2k/nogocomments/compare/v2.0.0...v3.0.0
[2.0.0]: https://github.com/pierow2k/nogocomments/compare/v1.0.3...v2.0.0
[1.0.3]: https://github.com/pierow2k/nogocomments/compare/v1.0.0...v1.0.3
[1.0.0]: https://github.com/pierow2k/nogocomments/releases/tag/v1.0.0
