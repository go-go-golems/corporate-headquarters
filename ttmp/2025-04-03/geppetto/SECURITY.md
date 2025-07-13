# Security Policy

## Reporting a Vulnerability

The go-go-golems team takes the security of our code seriously. We appreciate your efforts to responsibly disclose your findings and will make every effort to acknowledge your contributions.

To report a security vulnerability, please email [security@example.com](mailto:security@example.com) rather than using the public issue tracker. Include as much detail as possible to help us understand and reproduce the issue.

## Security Measures

This repository implements several security measures:

- Dependency scanning with Snyk, Dependabot, govulncheck, and nancy
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
