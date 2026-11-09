# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- TVBox web functionality integration
- Configuration management
- REST API endpoints
- WebSocket support
- Docker support

## [0.0.4] - 2026-11-09

### Added
- Port configuration support during service installation
- Config struct persisted in signed binary via ukey
- User prompted for port number (default: 8888) at install time
- Service reads port from signed config on startup

### Changed
- Refactored service.go for cleaner separation of concerns
- Extracted config logic into config.go

## [0.0.3] - 2026-11-09

### Fixed
- Fix service installation error by registering ukey buffer
- Import go-service buffer package to enable service signing
- Resolve "占位缓冲区未注册" error during service installation

### Technical Details
- Added import for `github.com/xxl6097/go-service/assets/buffer`
- This enables the ukey registration required for service signing
- Fixes installation failures on Linux servers

## [0.0.2] - 2026-11-09

### Added
- Release verification documentation
- Complete CI/CD validation report
- Verification workflow for release process

### Changed
- Improved project documentation structure
- Enhanced release verification process

### Documentation
- Added VERIFICATION_COMPLETE.md with comprehensive release validation
- Updated README with complete CI/CD workflow information
- Enhanced WORKFLOWS.md with detailed usage instructions
- Added CONTRIBUTING.md for community contributions

## [0.0.1] - 2026-11-09

### Added
- Initial project structure with modular design
- Service framework based on go-service
- Service lifecycle management (start, stop, shutdown)
- Logging system with zap logger integration
- Version information embedding
- Cross-platform support (Linux, macOS, Windows)
- Multi-architecture support (AMD64, ARM64)
- GitHub Actions CI/CD pipelines:
  - Automated build and test on push/PR
  - Automated multi-platform release
  - PR checks with linting and testing
  - Manual build workflow
- golangci-lint configuration
- Comprehensive documentation:
  - README with badges and quick start guide
  - Workflow usage guide
  - Contributing guidelines
- Makefile for local development

### Dependencies
- github.com/kardianos/service v1.2.4
- github.com/xxl6097/glog v0.1.96
- github.com/xxl6097/go-service v0.7.87
- go.uber.org/zap v1.28.0

### Notes
- This is the initial framework release
- TVBox web functionality is planned for future releases
- Service implements full IService interface from go-service

[Unreleased]: https://github.com/xxl6097/TVWeb/compare/v0.0.4...HEAD
[0.0.4]: https://github.com/xxl6097/TVWeb/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/xxl6097/TVWeb/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/xxl6097/TVWeb/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/xxl6097/TVWeb/releases/tag/v0.0.1
