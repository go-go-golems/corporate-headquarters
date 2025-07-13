# GitHub Actions Workflows for Supply Chain Security

Based on our research into supply chain security best practices for Go libraries with SOC2 compliance requirements, I'll now design GitHub Actions workflows to implement these security measures for the go-go-golems/geppetto repository.

## Key Security Requirements

From our research, we've identified these key security requirements:

1. **Dependency Management**: Track and scan dependencies for vulnerabilities
2. **Code Scanning**: Identify security vulnerabilities in code
3. **Secret Scanning**: Prevent secrets from being committed
4. **Automated Testing**: Ensure code quality and security
5. **Access Control**: Enforce proper authentication and authorization
6. **Change Management**: Ensure proper review and approval of changes
7. **Artifact Integrity**: Verify the integrity of build artifacts

## Workflow 1: Dependency Scanning and Management

This workflow will handle dependency scanning using Dependabot and Snyk.

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

## Workflow 2: Code Scanning with CodeQL

This workflow will perform static code analysis to identify security vulnerabilities.

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

## Workflow 3: Secret Scanning

This workflow will ensure no secrets are committed to the repository.

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

## Workflow 4: Comprehensive Testing

This workflow will run tests, linting, and other quality checks.

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

## Workflow 5: SLSA Build Verification

This workflow will implement SLSA (Supply chain Levels for Software Artifacts) level 2 compliance.

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

## Workflow 6: Release with Signing

This workflow will handle secure releases with artifact signing.

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

## Repository Settings and Branch Protection

In addition to workflows, implement these repository settings:

1. **Branch Protection Rules for `main` branch**:
   - Require pull request reviews before merging
   - Require status checks to pass before merging
   - Require signed commits
   - Do not allow bypassing the above settings

2. **Security Policy**:
   - Create a SECURITY.md file with vulnerability reporting instructions

3. **Dependabot Configuration**:
   ```yaml
   # .github/dependabot.yml
   version: 2
   updates:
     - package-ecosystem: "gomod"
       directory: "/"
       schedule:
         interval: "weekly"
       open-pull-requests-limit: 10
       target-branch: "develop"
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

These workflows and settings provide a comprehensive supply chain security setup for the go-go-golems/geppetto repository, addressing the key security requirements identified in our research and meeting SOC2 compliance requirements.
