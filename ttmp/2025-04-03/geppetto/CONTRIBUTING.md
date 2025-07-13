# Contributing to go-go-golems/geppetto

Thank you for your interest in contributing to our project! This document provides guidelines for contributing with a focus on security.

## Security Requirements

All contributions must adhere to these security requirements:

1. **Dependency Management**:
   - Do not add dependencies with known vulnerabilities
   - Use the minimum required version of dependencies
   - Document why new dependencies are needed

2. **Code Security**:
   - Follow Go security best practices
   - Avoid common vulnerabilities (SQL injection, XSS, etc.)
   - Handle errors appropriately

3. **Secret Management**:
   - Never commit secrets or credentials
   - Use environment variables for configuration
   - Follow the principle of least privilege

4. **Testing**:
   - Write tests for all new functionality
   - Include security-focused tests where appropriate

## Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests locally
5. Submit a pull request
6. Address review comments

All pull requests will undergo automated security checks. PRs that fail these checks will not be merged until the issues are resolved.

## Commit Signing

All commits must be signed with a GPG key. See GitHub's documentation on [signing commits](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits).

## Code of Conduct

Please follow our Code of Conduct in all your interactions with the project.
