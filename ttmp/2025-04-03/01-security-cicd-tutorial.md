# Setting Up Comprehensive Security and CI/CD for Go Repositories with GitHub Actions

This tutorial guides you through implementing a robust security and CI/CD setup for your Go repositories using GitHub Actions. We'll focus on implementing dependabot for automated dependency updates, along with various security scanning tools to ensure your codebase remains secure and up-to-date with minimal manual intervention.

## Why This Matters

Modern software development relies heavily on open-source dependencies. While these dependencies accelerate development, they also introduce potential security vulnerabilities. A comprehensive security setup helps detect:

- Outdated dependencies with known vulnerabilities
- Security issues in your codebase
- Accidental leakage of secrets and credentials
- Compliance issues for regulated industries

By automating these checks, you reduce the risk of security incidents while freeing up developer time to focus on building features.

## Prerequisites

- A GitHub repository with Go code
- Admin access to the repository
- Basic understanding of GitHub Actions

## Implementation Steps

We'll implement the following components:

1. **Dependabot Configuration** - Automated dependency updates for Go modules and GitHub Actions
2. **Dependency Scanning** - Detect vulnerabilities in dependencies
3. **Code Scanning** - Static analysis to find code-level security issues
4. **Secret Scanning** - Detect accidentally committed secrets
5. **Standard CI/CD** - Linting, testing, and release workflows

Let's get started!

## 1. Setting Up Dependabot

Dependabot automatically creates pull requests to update your dependencies. It supports both direct dependencies (Go modules) and GitHub Actions used in your workflows.

Create a file at `.github/dependabot.yml` with the following content:

```yaml
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

This configuration:
- Scans both Go modules and GitHub Actions weekly
- Adds "dependencies" and "security" labels to PRs for easy filtering
- Ignores patch updates for Go modules to reduce noise
- Limits to 10 open PRs at a time to prevent overwhelming the team

If your repository has multiple Go modules in different directories, you can add additional entries with different `directory` values.

## 2. Setting Up Dependency Scanning

Dependency scanning helps identify vulnerabilities in your dependencies. We'll set up multiple tools for comprehensive coverage:

Create a file at `.github/workflows/dependency-scanning.yml`:

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
        uses: actions/dependency-review-action@v4
        with:
          fail-on-severity: high

  govulncheck:
    name: Go Vulnerability Check
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Install govulncheck
        run: go install golang.org/x/vuln/cmd/govulncheck@latest
      
      - name: Run govulncheck
        run: govulncheck ./...
        
  nancy:
    name: Nancy Vulnerability Scan
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Install Nancy
        run: go install github.com/sonatype-nexus-community/nancy@latest
      
      - name: Run Nancy
        run: go list -json -deps ./... | nancy sleuth

  gosec:
    name: GoSec Security Scan
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Run Gosec Security Scanner
        uses: securego/gosec@master
        with:
          args: -exclude=G101,G304,G301,G306 -exclude-dir=.history ./...
```

This workflow includes:
- **Dependency Review** - Scans for vulnerabilities during pull requests
- **govulncheck** - Go-specific vulnerability scanner from the Go team
- **Nancy** - Alternative dependency scanner with different vulnerability database
- **Gosec** - Security linter for Go code that finds common security issues

## 3. Setting Up Code Scanning with CodeQL

CodeQL is GitHub's semantic code analysis engine that finds vulnerabilities and errors in your code.

Create a file at `.github/workflows/codeql-analysis.yml`:

```yaml
name: "CodeQL Analysis"

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '0 0 * * 0'  # Run weekly on Sunday at midnight

jobs:
  analyze:
    name: Analyze
    runs-on: ubuntu-latest
    permissions:
      security-events: write
      
    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Initialize CodeQL
      uses: github/codeql-action/init@v3
      with:
        languages: go

    - name: Perform CodeQL Analysis
      uses: github/codeql-action/analyze@v3
```

CodeQL performs deep semantic analysis of your code to find security vulnerabilities and coding errors that might not be detected by other static analysis tools.

## 4. Setting Up Secret Scanning

Secret scanning detects accidentally committed secrets, such as API keys or credentials.

Create a file at `.github/workflows/secret-scanning.yml`:

```yaml
name: Secret Scanning

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  trufflehog:
    name: TruffleHog Secret Scan
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
```

TruffleHog scans git history to detect secrets that may have been committed accidentally. It's particularly useful for finding credentials that might have been committed and later removed but still exist in the git history.

## 5. Setting Up Standard CI/CD Workflows

In addition to security scanning, we'll set up standard CI/CD workflows for Go projects:

### 5.1 Linting

Create a file at `.github/workflows/lint.yml`:

```yaml
name: golangci-lint
on:
  push:
    tags:
      - v*
    branches:
      - main
  pull_request:
permissions:
  contents: read

jobs:
  golangci:
    permissions:
      contents: read
      pull-requests: read
    name: lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - run: git fetch --force --tags
      - uses: actions/setup-go@v5
        with:
          go-version: '>=1.19.5'
          cache: true
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v7
        with:
          version: v2.0.2
          args: --timeout=5m
```

### 5.1.1 Migrating to golangci-lint v2

As of 2024, golangci-lint released version 2.0, which includes significant changes to the configuration format. If you're upgrading from v1, you'll need to update your `.golangci.yml` configuration file.

Here's an example of migrating from v1 to v2:

**Old config (v1):**
```yaml
linters:
  disable-all: false
  enable:
    - errcheck
    - govet
    - ineffassign
    - staticcheck
    - typecheck
    - unused
    - gofmt
  fast: false
issues:
  exclude:
    - 'SA1019: cli.CreateProcessorLegacy'
```

**New config (v2):**
```yaml
# Defines the configuration version.
# The only possible value is "2".
version: "2"

# Linters configuration
linters:
  # Default set of linters.
  default: none
  # Enable specific linters
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
  # Exclusions configuration
  exclusions:
    rules:
      - linters:
          - staticcheck
        text: 'SA1019: cli.CreateProcessorLegacy'

# Formatters configuration
formatters:
  enable:
    - gofmt
```

Key migration changes:
1. Add `version: "2"` to indicate v2 configuration format
2. Replace `disable-all: false/true` with `default: all/none`
3. The `typecheck` linter is now integrated into Go's type checking system
4. Move formatters (`gofmt`, `gofumpt`, `goimports`, `gci`) from `linters.enable` to a separate `formatters.enable` section
5. Replace `issues.exclude` with more structured `linters.exclusions.rules`
6. The `fast` option is removed (use `--fast-only` flag or `default: fast` instead)

### 5.1.2 Setting Up Docker-based Linting

To ensure consistent linting across different development environments, add Docker-based linting to your Makefile:

```makefile
docker-lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v2.0.2 golangci-lint run -v

lint:
	golangci-lint run -v

lintmax:
	golangci-lint run -v --max-same-issues=100
```

### 5.1.3 Additional Security Scanning Targets

Add these security scanning targets to your Makefile:

```makefile
gosec:
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec -exclude=G101,G304,G301,G306 -exclude-dir=.history ./...

govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...
```

### 5.1.4 Configuring GoSec

Create a `.gosec.config.json` file in your project root to configure GoSec behavior:

```json
{
  "exclude": {
    "paths": [".history/"]
  }
}
```

This configuration:
- Excludes the `.history/` directory from security scanning
- Helps prevent false positives from temporary or generated files
- Can be extended with additional paths or rules as needed

### 5.1.5 Git Hooks with Lefthook

[Lefthook](https://github.com/evilmartians/lefthook) provides a fast and powerful Git hooks manager. Create a `lefthook.yml` file in your project root:

```yaml
pre-commit:
  commands:
    lint:
      glob: '*.go'
      run: make lintmax gosec govulncheck
    test:
      glob: '*.go'
      run: make test
  parallel: true

pre-push:
  commands:
    release:
      run: make goreleaser
    lint:
      run: make lintmax gosec govulncheck
    test:
      run: make test
  parallel: true
```

This configuration:
- Runs linting, security checks, and tests before each commit
- Performs additional checks including release validation before pushing
- Executes tasks in parallel for better performance
- Only triggers on Go file changes for relevant hooks

To install Lefthook:

```bash
# Using Go
go install github.com/evilmartians/lefthook@latest

# Or using Homebrew
brew install lefthook

# Initialize in your repository
lefthook install
```

## 6. Automating Security Checks

With these configurations in place, security checks are automated at multiple levels:

1. **Local Development**
   - Pre-commit hooks run linting and security checks
   - Pre-push hooks ensure all tests pass
   - Make targets available for manual security scanning

2. **CI/CD Pipeline**
   - GitHub Actions run comprehensive security checks
   - Dependabot keeps dependencies updated
   - CodeQL performs deep code analysis

3. **Release Process**
   - Security checks must pass before releases
   - Automated vulnerability scanning of dependencies
   - Signed releases with proper versioning

## 7. Enabling GitHub Security Features

For these workflows to be fully effective, you need to enable several security features in your GitHub repository settings:

1. Go to your repository on GitHub
2. Click on "Settings"
3. In the left sidebar, click on "Code security and analysis"
4. Enable the following features:
   - Dependency graph
   - Dependabot alerts
   - Dependabot security updates
   - Code scanning
   - Secret scanning

## 8. Keeping Actions Up-to-Date

GitHub Actions are frequently updated with new features, bug fixes, and security improvements. To ensure you're using the latest versions, periodically check for updates to the actions you're using.

The dependabot configuration we set up earlier will automatically create pull requests when new versions of actions are available.

## Best Practices for Maintenance

1. **Review Dependabot PRs regularly** - Don't let them pile up
2. **Address security alerts promptly** - Especially those with high severity
3. **Keep Go and action versions updated** - Newer versions often include security fixes
4. **Monitor workflow runs** - Check for failures and address issues
5. **Periodically review security settings** - New security features are added regularly

## Conclusion

With this setup, your Go repository is now equipped with comprehensive security scanning and CI/CD automation. These workflows will help detect potential security issues early, keep dependencies up-to-date, and ensure code quality through automated testing and linting.

By automating these checks, you reduce the maintenance burden on your team while improving the security posture of your project. You're now following industry best practices for secure Go development with GitHub Actions.

## Further Reading

- [GitHub Dependabot Documentation](https://docs.github.com/en/code-security/dependabot/dependabot-version-updates)
- [Go Security Best Practices](https://go.dev/security/best-practices)
- [GitHub Advanced Security](https://docs.github.com/en/github/getting-started-with-github/learning-about-github/about-github-advanced-security)
- [GoReleaser Documentation](https://goreleaser.com/) 