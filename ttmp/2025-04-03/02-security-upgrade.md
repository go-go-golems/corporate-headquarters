# Upgrading Your Go Repository Security and CI/CD Setup

This guide complements the `01-security-cicd-tutorial.md` and provides step-by-step instructions on how to apply the recommended security and CI/CD enhancements to an existing Go project based on the provided configuration changes.

## Introduction

This upgrade integrates automated dependency updates (Dependabot), comprehensive dependency and code scanning (Dependency Review, govulncheck, Nancy, Gosec, CodeQL), secret scanning (TruffleHog), and updated CI/CD practices (golangci-lint v2, Lefthook integration) into your Go repository.

## Prerequisites

- An existing Go project managed with Git and hosted on GitHub.
- Familiarity with GitHub Actions and Makefiles.
- Admin access to the repository to configure settings and secrets if needed.

## Upgrade Steps

Follow these steps to apply the changes:

### 1. Add Dependabot Configuration

Dependabot automates dependency updates for Go modules and GitHub Actions.

Create the file `.github/dependabot.yml` with the following content:

```yaml
# .github/dependabot.yml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/" # Adjust if your go.mod is not in the root
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
* **Note:** Adjust the `directory` for `gomod` if your `go.mod` file is not in the repository root.

### 2. Add Dependency Scanning Workflow

This workflow combines multiple tools to scan for vulnerabilities in your dependencies.

Create the file `.github/workflows/dependency-scanning.yml`:

```yaml
# .github/workflows/dependency-scanning.yml
name: Dependency Scanning

on:
  push:
    branches: [ main ] # Adjust if your main branch has a different name
  pull_request:
    branches: [ main ] # Adjust if your main branch has a different name
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
          go-version: '1.24' # Ensure this matches your project's Go version

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
          go-version: '1.24' # Ensure this matches your project's Go version

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
          go-version: '1.24' # Ensure this matches your project's Go version

      - name: Run Gosec Security Scanner
        uses: securego/gosec@master
        with:
          # Adjust exclusions based on your project's needs
          args: -exclude=G101,G304,G301,G306,G204 -exclude-dir=.history ./...
```
* **Note:** Update the `go-version` to match your project's requirement. Review and adjust the `gosec` arguments (`-exclude`, `-exclude-dir`) as needed for your project.

### 3. Add CodeQL Analysis Workflow

CodeQL performs deep static analysis to find potential vulnerabilities.

Create the file `.github/workflows/codeql-analysis.yml`:

```yaml
# .github/workflows/codeql-analysis.yml
name: "CodeQL Analysis"

on:
  push:
    branches: [ main ] # Adjust if your main branch has a different name
  pull_request:
    branches: [ main ] # Adjust if your main branch has a different name
  schedule:
    - cron: '0 0 * * 0'  # Run weekly on Sunday at midnight

jobs:
  analyze:
    name: Analyze
    runs-on: ubuntu-latest
    permissions:
      security-events: write # Required for CodeQL to report findings

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

### 4. Add Secret Scanning Workflow

This workflow uses TruffleHog to scan for accidentally committed secrets in your Git history.

Create the file `.github/workflows/secret-scanning.yml`:

```yaml
# .github/workflows/secret-scanning.yml
name: Secret Scanning

on:
  push:
    branches: [ main ] # Adjust if your main branch has a different name
  pull_request:
    branches: [ main ] # Adjust if your main branch has a different name

jobs:
  trufflehog:
    name: TruffleHog Secret Scan
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0 # Required for TruffleHog to scan history

      - name: TruffleHog OSS
        uses: trufflesecurity/trufflehog@main
        with:
          path: ./
          base: ${{ github.event.repository.default_branch }}
          head: HEAD
```

### 5. Update Lint Workflow

This updates the linting workflow to use `golangci-lint-action@v7` and specifies version `v2.0.2`.

Modify `.github/workflows/lint.yml`:

```diff
--- a/.github/workflows/lint.yml
+++ b/.github/workflows/lint.yml
@@ -27,6 +27,7 @@ jobs:
           go-version: '>=1.19.5' # Or your required Go version range
           cache: true
       - name: golangci-lint
-        uses: golangci/golangci-lint-action@v6 # Or your previous version
+        uses: golangci/golangci-lint-action@v7
         with:
+          version: v2.0.2
           args: --timeout=5m # Keep or adjust args as needed
```
* **Important:** If you were using `golangci-lint` v1, you **must** migrate your `.golangci.yml` configuration to the v2 format. Refer to section `5.1.1 Migrating to golangci-lint v2` in `01-security-cicd-tutorial.md` for details.

### 5.1 Migrate .golangci.yml to v2 Format

Since you're updating to golangci-lint v2, you need to migrate your configuration file. Here's an example migration:

```diff
--- a/.golangci.yml
+++ b/.golangci.yml
@@ -1,22 +1,26 @@
+version: "2"
+
 linters:
-  disable-all: false
+  default: none
   enable:
     # defaults
     - errcheck
-    - gosimple
     - govet
     - ineffassign
     - staticcheck
-    - typecheck
     - unused
     # stuff I'm adding
     - exhaustive
     #    - gochecknoglobals
     #    - gochecknoinits
-    - gofmt
     - nonamedreturns
     - predeclared
-  fast: false
-issues:
-  exclude:
-    - 'SA1019: cli.CreateProcessorLegacy'
+  exclusions:
+    rules:
+      - linters:
+          - staticcheck
+        text: 'SA1019: cli.CreateProcessorLegacy'
+
+formatters:
+  enable:
+    - gofmt
```

Key changes in this migration:
1. Add `version: "2"` at the top to specify v2 configuration format
2. Replace `disable-all: false` with `default: none` 
3. Remove `gosimple` (now included in `staticcheck` in v2)
4. Remove `typecheck` (now integrated into Go's type checking system)
5. Move `gofmt` from `linters.enable` to a new `formatters.enable` section
6. Replace `issues.exclude` with structured `linters.exclusions.rules`
7. Remove the `fast` option (no longer used in v2)

Make sure to test the new configuration locally before committing. The simplest way to test is to run:
```bash
golangci-lint run -v
```

### 6. Update Release Workflow

This updates the versions of actions used for GPG key import and GoReleaser.

Modify `.github/workflows/release.yaml` (or `.yml`):

```diff
--- a/.github/workflows/release.yaml
+++ b/.github/workflows/release.yaml
@@ -28,13 +28,13 @@ jobs:

       - name: Import GPG key
         id: import_gpg
-        uses: crazy-max/ghaction-import-gpg@v5 # Or your previous version
+        uses: crazy-max/ghaction-import-gpg@v6
         with:
           gpg_private_key: ${{ secrets.GO_GO_GOLEMS_SIGN_KEY }}
           passphrase: ${{ secrets.GO_GO_GOLEMS_SIGN_PASSPHRASE }}
           fingerprint: "6EBE1DF0BDF48A1BBA381B5B79983EF218C6ED7E" # Your GPG key fingerprint

-      - uses: goreleaser/goreleaser-action@v4 # Or your previous version
+      - uses: goreleaser/goreleaser-action@v6
         with:
           distribution: goreleaser
           version: latest
```
* **Note:** Ensure the GPG secrets (`GO_GO_GOLEMS_SIGN_KEY`, `GO_GO_GOLEMS_SIGN_PASSPHRASE`) and `fingerprint` are correctly configured for your repository if you use signed releases.

### 6.1 Update GoReleaser Configuration

Along with updating the GoReleaser GitHub Action, you should also update your GoReleaser configuration file to ensure compatibility with the latest version and to improve build architecture coverage.

Modify `.goreleaser.yaml`:

```diff
--- a/.goreleaser.yaml
+++ b/.goreleaser.yaml
@@ -17,6 +17,9 @@ builds:
 # I am not able to test windows at the time
 #      - windows
       - darwin
+    goarch:
+      - amd64
+      - arm64
 checksum:
   name_template: 'checksums.txt'
 snapshot:
@@ -47,7 +50,7 @@ nfpms:
       - rpm

     # Version Release.
-    release: 1
+    release: "1"

     # Section.
     section: default
```

This update:
* Explicitly defines target architectures (`amd64` and `arm64`) to ensure proper cross-platform support, particularly for modern Mac systems with Apple Silicon
* Fixes the `release` value format by making it a string (`"1"`) instead of an integer (`1`), which is required by newer GoReleaser versions for proper parsing

Additionally, review your complete `.goreleaser.yaml` to ensure it follows current best practices:
* Consider adding `universal_binaries` section for macOS if appropriate
* Verify signing configuration if you use signed releases
* Check if any deprecated fields or features are used

### 7. Update Makefile

Add new targets for enhanced linting and security scanning, update the Docker lint image version, and potentially add dependency update targets like `bump-glazed`.

Apply the following changes to your `Makefile`:

```diff
--- a/Makefile
+++ b/Makefile
@@ -1,32 +1,32 @@
-.PHONY: gifs
+.PHONY: all test build lint lintmax docker-lint gosec govulncheck goreleaser tag-major tag-minor tag-patch release bump-glazed install # Add new targets

-VERSION ?= $(shell git describe --tags --always --dirty)
-COMMIT ?= $(shell git rev-parse --short HEAD)
-DIRTY ?= $(shell git diff --quiet || echo "dirty")
+all: test build # Adjust 'all' target as needed

-LDFLAGS=-ldflags "-X main.version=$(VERSION)-$(COMMIT)-$(DIRTY)"
-
-all: gifs
+VERSION=vX.Y.Z # Update if you manage version here, or remove if using GoReleaser exclusively

-TAPES=$(shell ls doc/vhs/*tape)
-gifs: $(TAPES)
-	for i in $(TAPES); do vhs < $$i; done
-
 docker-lint:
-	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run -v
+	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v2.0.2 golangci-lint run -v

 lint:
-	golangci-lint run -v --enable=exhaustive
+	golangci-lint run -v # Basic lint run

+lintmax: # New target for more comprehensive local linting
+	golangci-lint run -v --max-same-issues=100
+
+gosec: # New target for GoSec scan
+	go install github.com/securego/gosec/v2/cmd/gosec@latest
+	# Adjust exclusions as needed, mirroring the workflow
+	gosec -exclude=G101,G304,G301,G306,G204 -exclude-dir=.history ./...
+
+govulncheck: # New target for govulncheck scan
+	go install golang.org/x/vuln/cmd/govulncheck@latest
+	govulncheck ./...

 test:
 	go test ./...

 build:
 	go generate ./...
-	go build $(LDFLAGS) ./...
-
-bench:
-	go test -bench=./... -benchmem
+	go build ./... # Remove LDFLAGS if not needed or handled by GoReleaser

 goreleaser:
 	goreleaser release --skip=sign --snapshot --clean
@@ -44,11 +44,11 @@ release:
 	git push --tags
 	GOPROXY=proxy.golang.org go list -m github.com/go-go-golems/glazed@$(shell svu current)
 
-exhaustive: # Remove old exhaustive target if present
-	golangci-lint run -v --enable=exhaustive
-
-GLAZE_BINARY=$(shell which glaze)
+bump-glazed: # Example dependency bump target, adapt as needed
+	go get github.com/go-go-golems/glazed@latest
+	go get github.com/go-go-golems/clay@latest # Add other dependencies if needed
+	go mod tidy

 install:
-	go build $(LDFLAGS) -o ./dist/glaze ./cmd/glaze && \
-		cp ./dist/glaze $(GLAZE_BINARY)
+	go build -o ./dist/yourbinary ./cmd/yourbinary && \ # Update binary name
+		cp ./dist/yourbinary $(shell which yourbinary) # Update binary name

### 7.1 Add CodeQL Local Analysis Target

Adding a target to run CodeQL analysis locally allows you to catch security vulnerabilities before pushing code to GitHub. This complements the GitHub Actions workflow and enables offline security analysis.

Add the following to your `Makefile`:

```makefile
# Path to CodeQL CLI - adjust based on installation location
CODEQL_PATH ?= $(shell which codeql)
# Path to CodeQL queries - adjust based on where you cloned the repository
CODEQL_QUERIES ?= $(HOME)/codeql-go/ql/src/go

# Create CodeQL database and run analysis
codeql-local:
	@if [ -z "$(CODEQL_PATH)" ]; then echo "CodeQL CLI not found. Install from https://github.com/github/codeql-cli-binaries/releases"; exit 1; fi
	@if [ ! -d "$(CODEQL_QUERIES)" ]; then echo "CodeQL queries not found. Clone from https://github.com/github/codeql-go"; exit 1; fi
	$(CODEQL_PATH) database create --language=go --source-root=. ./codeql-db
	$(CODEQL_PATH) database analyze ./codeql-db $(CODEQL_QUERIES)/Security --format=sarif-latest --output=codeql-results.sarif
	@echo "Results saved to codeql-results.sarif"
```

#### Installing CodeQL CLI

Before you can use the `codeql-local` target, you need to set up the CodeQL CLI:

1. Download the latest CodeQL CLI from [GitHub's releases](https://github.com/github/codeql-cli-binaries/releases)
2. Extract the archive to a folder (e.g., `~/codeql-cli`)
3. Add the CLI to your PATH:
   ```bash
   echo 'export PATH="$HOME/codeql-cli:$PATH"' >> ~/.bashrc
   source ~/.bashrc
   ```
4. Verify the installation:
   ```bash
   codeql --version
   ```

#### Getting the CodeQL Query Libraries

You also need the CodeQL query libraries for Go:

```bash
git clone https://github.com/github/codeql-go.git ~/codeql-go
```

#### Running the Analysis

Once set up, you can run:

```bash
make codeql-local
```

This will:
1. Create a CodeQL database by analyzing your Go code
2. Run security queries against the database
3. Output the results in SARIF format to `codeql-results.sarif`

You can view the results with any SARIF viewer or with tools like Visual Studio Code (using the SARIF extension).

### 8. Update Lefthook Configuration

Update Git hooks managed by Lefthook to run the new, more comprehensive checks.

Modify `lefthook.yml`:

```diff
--- a/lefthook.yml
+++ b/lefthook.yml
@@ -2,13 +2,10 @@ pre-commit:
   commands:
     lint:
       glob: '*.go'
-      run: make lint # Or your previous lint command
+      run: make lintmax gosec govulncheck # Run new comprehensive checks
     test:
       glob: '*.go'
       run: make test
-    exhaustive: # Remove old targets if present
-      glob: '*.go'
-      run: make exhaustive
   parallel: true

 pre-push:
@@ -16,7 +13,7 @@ pre-push:
     release:
       run: make goreleaser
     lint:
-      run: make lint
+      run: make lintmax gosec govulncheck # Run new comprehensive checks
     test:
       run: make test
   parallel: true

```
* **Note:** If you aren't using Lefthook, you can install it and initialize it (see section `5.1.5` in the main tutorial) or adapt the commands for your existing Git hook setup. Run `lefthook install` after modifying the `lefthook.yml` file to update the hooks.

### 9. Run Initial Scans and Fix Issues

After applying the configuration changes, run the new Make targets locally to identify and fix any existing issues.

You can run all scans at once to get a comprehensive view of issues that need to be addressed:

```bash
make lintmax; make gosec; make govulncheck; make test
```

Then address each category of issues one by one:

1.  **Linting Issues:**
    * Fix any linting errors reported by `golangci-lint`
    * Adjust your code to comply with the enabled linters
    * If needed, fine-tune your `.golangci.yml` configuration to customize rules or exclude specific violations

2.  **Security Issues (GoSec):**
    * Review security vulnerabilities identified by GoSec
    * For each issue:
        * If it's a genuine vulnerability, fix the code
        * If it's a false positive or acceptable risk, add a `#nosec GXXX` comment (where `GXXX` is the rule ID) with a justification
        * Example: `_ = os.Setenv("DANGEROUS", "true") // #nosec G104 -- Safe because we control the value`
    * Consider adding path exclusions to a `.gosec.config.json` file for directories that shouldn't be scanned

3.  **Dependency Vulnerabilities (govulncheck):**
    * Review any reported vulnerabilities in your dependencies
    * Update vulnerable modules using `go get module@latest` or a specific safe version
    * Run `go mod tidy` after updates
    * Verify fixes with another run of `make govulncheck`

### 10. Enable GitHub Security Features

Ensure the necessary security features are enabled in your GitHub repository settings.

1.  Go to your repository on GitHub -> Settings -> Code security and analysis.
2.  Enable:
    *   Dependency graph
    *   Dependabot alerts
    *   Dependabot security updates
    *   Code scanning alerts (using CodeQL and any other tools you've integrated)
    *   Secret scanning alerts

## Conclusion

By completing these steps, you have significantly upgraded your repository's security posture and CI/CD automation. Your project now benefits from automated dependency management, multiple layers of vulnerability and code scanning, secret detection, and streamlined local and remote checks. Remember to regularly review Dependabot PRs and address security alerts reported by the new workflows.