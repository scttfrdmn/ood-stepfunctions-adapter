# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-05-30

### Added
- Initial scaffold — OOD compute adapter for AWS Step Functions, translating Open OnDemand job submissions into Step Functions executions.
- CLI commands: `submit` (JSON job spec from stdin → start an execution, prints the execution ARN), `status <execution-arn>` (OOD-normalized status), `delete <execution-arn>` (stop a running execution, maps to OOD job cancel), and `info <execution-arn>` (full execution detail as JSON).
- Unit and substrate integration tests.

### Changed
- Upgraded to substrate v0.45.2 and removed the local `replace` directive; earlier upgraded to v0.44.3 (fixes substrate #242 Step Functions routing bug).
