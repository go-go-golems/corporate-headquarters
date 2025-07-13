# Implementation Guide: Supply Chain Security Setup for go-go-golems/geppetto

This implementation guide provides step-by-step instructions for setting up a comprehensive supply chain security system for the go-go-golems/geppetto Go library repository. The setup is designed to meet SOC2 compliance requirements while leveraging GitHub Actions and complementary security tools.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Repository Configuration](#repository-configuration)
3. [GitHub Actions Workflows Implementation](#github-actions-workflows-implementation)
4. [Additional Security Tools Integration](#additional-security-tools-integration)
5. [Security Policy and Documentation](#security-policy-and-documentation)
6. [Validation and Monitoring](#validation-and-monitoring)
7. [Maintenance Procedures](#maintenance-procedures)

## Prerequisites

Before implementing this security setup, ensure you have:

- Administrator access to the go-go-golems/geppetto GitHub repository
- A Snyk account with an API token (already in use)
- GPG key for signing commits and releases (if not already set up)
- GitHub Actions enabled for the repository

## Repository Configuration

### 1. Branch Protection Rules

Set up branch protection rules for the main branch:

1. Navigate to the repository on GitHub
2. Go to Settings > Branches
3. Click "Add rule" next to Branch protection rules
4. Enter "main" as the branch name pattern
5. Configure the following settings:
   - ✅ Require pull request reviews before merging
   - ✅ Require status checks to pass before merging
   - ✅ Require signed commits
   - ✅ Include administrators
   - ✅ Restrict who can push to matching branches
6. Click "Create" or "Save changes"

### 2. Repository Settings

Configure general repository security settings:

1. Go to Settings > Security & analysis
2. Enable the following features:
   - ✅ Dependency graph
   - ✅ Dependabot alerts
   - ✅ Dependabot security updates
   - ✅ Code scanning
   - ✅ Secret scanning

### 3. Required Secrets Setup

Add the following secrets to your repository:

1. Go to Settings > Secrets and variables > Actions
2. Add the following repository secrets:
   - `SNYK_TOKEN`: Your Snyk API token
   - `GPG_PRIVATE_KEY`: Your GPG private key for signing
   - `PASSPHRASE`: The passphrase for your GPG key

## GitHub Actions Workflows Implementation

Create the following workflow files in the `.github/workflows/` directory:

### 1. Dependency Scanning Workflow

Create a file named `.github/workflows/dependency-scanning.yml`:

```yaml
# .github/workflows/dependency-scanning.yml
name: Dependency Scanning

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '0 0 * * 0'  # Run weekly on Sunday at midnight

jobs:
  dependency-review:
    name: Dependency Review
    runs-on: ubuntu-latest
    if: github.event_name == 'pull_request'
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Dependency Review
        uses: actions/dependency-review-action@v3
        with:
          fail-on-severity: high

  snyk-scan:
    name: Snyk Vulnerability Scan
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run Snyk to check for vulnerabilities
        uses: snyk/actions/golang@master
        env:
          SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
        with:
          args: --severity-threshold=high

  govulncheck:
    name: Go Vulnerability Check
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install govulncheck
        run: go install golang.org/x/vuln/cmd/govulncheck@latest
      
      - name: Run govulncheck
        run: govulncheck ./...
```

### 2. Code Scanning Workflow

Create a file named `.github/workflows/codeql-analysis.yml`:

```yaml
# .github/workflows/codeql-analysis.yml
name: "CodeQL Analysis"

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '0 0 * * 1'  # Run weekly on Monday at midnight

jobs:
  analyze:
    name: Analyze
    runs-on: ubuntu-latest
    permissions:
      security-events: write
      actions: read
      contents: read

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Initialize CodeQL
        uses: github/codeql-action/init@v2
        with:
          languages: go

      - name: Perform CodeQL Analysis
        uses: github/codeql-action/analyze@v2
```

### 3. Secret Scanning Workflow

Create a file named `.github/workflows/secret-scanning.yml`:

```yaml
# .github/workflows/secret-scanning.yml
name: Secret Scanning

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  secret-scan:
    name: Detect Secrets
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: TruffleHog OSS
        uses: trufflesecurity/trufflehog@main
        with:
          path: ./
          base: ${{ github.event.repository.default_branch }}
          head: HEAD
          extra_args: --debug --only-verified
```

### 4. Testing Workflow

Create a file named `.github/workflows/test.yml`:

```yaml
# .github/workflows/test.yml
name: Test and Lint

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Verify dependencies
        run: go mod verify
      
      - name: Build
        run: go build -v ./...
      
      - name: Run tests with coverage
        run: go test -race -coverprofile=coverage.txt -covermode=atomic ./...
      
      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.txt

  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
```

### 5. SLSA Build Verification Workflow

Create a file named `.github/workflows/slsa-build.yml`:

```yaml
# .github/workflows/slsa-build.yml
name: SLSA Go Build

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    name: SLSA Build
    runs-on: ubuntu-latest
    permissions:
      id-token: write
      contents: read
      actions: read
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      # For a real implementation, use the SLSA generator
      - name: SLSA Provenance Generator
        uses: slsa-framework/slsa-github-generator/.github/workflows/generator_go_slsa3.yml@v1.4.0
        with:
          go-version: '1.21'
          artifact-path: ./bin/
```

### 6. Release Workflow

Create a file named `.github/workflows/release.yml`:

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  goreleaser:
    name: Release
    runs-on: ubuntu-latest
    permissions:
      contents: write
      packages: write
      id-token: write
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Import GPG key
        id: import_gpg
        uses: crazy-max/ghaction-import-gpg@v6
        with:
          gpg_private_key: ${{ secrets.GPG_PRIVATE_KEY }}
          passphrase: ${{ secrets.PASSPHRASE }}
      
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GPG_FINGERPRINT: ${{ steps.import_gpg.outputs.fingerprint }}
```

### 7. Dependabot Configuration

Create a file named `.github/dependabot.yml`:

```yaml
# .github/dependabot.yml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
    labels:
      - "dependencies"
      - "security"
    ignore:
      - dependency-name: "*"
        update-types: ["version-update:semver-patch"]
  
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
    labels:
      - "dependencies"
      - "security"
```

## Additional Security Tools Integration

In addition to the GitHub Actions workflows, integrate these recommended security tools:

### 1. gosec - Go Security Checker

Add gosec to your testing workflow by updating the `.github/workflows/test.yml` file:

```yaml
# Add this job to the existing test.yml file
  gosec:
    name: Go Security Check
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Run Gosec Security Scanner
        uses: securego/gosec@master
        with:
          args: -exclude=G101,G304,G301,G306 -exclude-dir=.history ./...
```

### 2. SonarCloud Integration

Create a new file named `.github/workflows/sonarcloud.yml`:

```yaml
# .github/workflows/sonarcloud.yml
name: SonarCloud Analysis

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  sonarcloud:
    name: SonarCloud
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests with coverage
        run: go test -coverprofile=coverage.out ./...
      
      - name: SonarCloud Scan
        uses: SonarSource/sonarcloud-github-action@master
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
```

Note: You'll need to add a `SONAR_TOKEN` secret to your repository and create a `sonar-project.properties` file in your repository root.

### 3. OSSF Scorecard

Create a new file named `.github/workflows/scorecard.yml`:

```yaml
# .github/workflows/scorecard.yml
name: OSSF Scorecard

on:
  schedule:
    - cron: '0 0 * * 0'  # Run weekly on Sunday at midnight
  push:
    branches: [ main ]

jobs:
  scorecard:
    name: OSSF Scorecard Analysis
    runs-on: ubuntu-latest
    permissions:
      security-events: write
      id-token: write
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          persist-credentials: false
      
      - name: Run OSSF Scorecard
        uses: ossf/scorecard-action@v2.1.2
        with:
          results_file: results.sarif
          results_format: sarif
          publish_results: true
      
      - name: Upload SARIF results
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: results.sarif
```

### 4. nancy - Go Dependency Vulnerability Scanner

Add nancy to your dependency scanning workflow by updating the `.github/workflows/dependency-scanning.yml` file:

```yaml
# Add this job to the existing dependency-scanning.yml file
  nancy:
    name: Nancy Vulnerability Scan
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install Nancy
        run: go install github.com/sonatype-nexus-community/nancy@latest
      
      - name: Run Nancy
        run: go list -json -deps ./... | nancy sleuth
```

## Security Policy and Documentation

### 1. Create a Security Policy

Create a file named `SECURITY.md` in the repository root:

```markdown
# Security Policy

## Reporting a Vulnerability

The go-go-golems team takes the security of our code seriously. We appreciate your efforts to responsibly disclose your findings and will make every effort to acknowledge your contributions.

To report a security vulnerability, please email [security@example.com](mailto:security@example.com) rather than using the public issue tracker. Include as much detail as possible to help us understand and reproduce the issue.

## Security Measures

This repository implements several security measures:

- Dependency scanning with Snyk, Dependabot, and govulncheck
- Code scanning with CodeQL and gosec
- Secret scanning with TruffleHog
- SLSA build verification
- Signed commits and releases

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |
| < latest| :x:                |

## Security Update Process

When we receive a security report, we will:

1. Confirm the issue and determine its severity
2. Develop and test a fix
3. Prepare a security advisory
4. Release the fix and publish the advisory

## Disclosure Policy

When we receive a security bug report, we will:

- Confirm the problem and determine the affected versions
- Audit code to find any similar problems
- Prepare fixes for all supported versions
- Release new versions and notify users
```

### 2. Developer Documentation

Create a file named `CONTRIBUTING.md` with security guidelines for contributors:

```markdown
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

Please follow our [Code of Conduct](CODE_OF_CONDUCT.md) in all your interactions with the project.
```

## Validation and Monitoring

### 1. Initial Validation

After implementing the security setup:

1. Run all workflows manually to ensure they work correctly
2. Check for any false positives and adjust configurations as needed
3. Verify that branch protection rules are working as expected
4. Test the Dependabot configuration by introducing a vulnerable dependency

### 2. Ongoing Monitoring

Set up monitoring to ensure continued security:

1. Review security alerts regularly
2. Monitor workflow execution results
3. Periodically review access controls and permissions
4. Schedule quarterly security reviews

### 3. Compliance Validation

Fo
(Content truncated due to size limit. Use line ranges to read in chunks)