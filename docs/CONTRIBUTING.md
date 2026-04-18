# Contributing to ICD Codes API

Thank you for your interest in improving the ICD Codes API. We follow industry best practices for Go development and maintain a modular, domain-driven architecture.

## Development Workflow
1. Fork the repository.
2. Create a feature branch: `git checkout -b feat/your-feature`.
3. Ensure your code follows the established project structure:
   - `internal/` for private implementation logic.
   - `cmd/` for entry points.
4. Submit a Pull Request with a clear description of the problem solved.

## Commit Standards
We use [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` for new capabilities.
- `fix:` for bug fixes.
- `chore:` for maintenance or documentation.

## Architecture
The project is built using Domain-Driven Design (DDD). We maintain a strict separation between:
- **Handlers**: Web-specific logic (Gin).
- **Services**: Business logic.
- **Repositories**: Data access.
