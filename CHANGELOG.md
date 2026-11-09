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

[Unreleased]: https://github.com/xxl6097/TVWeb/compare/v0.0.1...HEAD
[0.0.1]: https://github.com/xxl6097/TVWeb/releases/tag/v0.0.1
