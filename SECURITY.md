# nogocomments Security Policy

This security policy applies to the nogocomments GitHub repository
and outlines the process for reporting security vulnerabilities, handling
security incidents, and ensuring the safety of our users. It adheres to
the principles of Coordinated Vulnerability Disclosure (CVD), fostering
collaboration between reporters and the nogocomments team to address
vulnerabilities responsibly and effectively.

## Overview

nogocomments is a command-line interface (CLI) tool designed with
security and reliability in mind. While its attack surface may be limited,
we remain vigilant in addressing potential vulnerabilities to protect our
users and maintain trust in the nogocomments ecosystem.

## Reporting Security Issues

If you discover a potential security issue in nogocomments, we ask
that you report it to us privately to allow us the opportunity to address
it before public disclosure.

### Reporting Methods

- **Email**: [security@daspyro.de](mailto:security@daspyro.de)
- **GitHub Security Advisories**: [GitHub Advisories Documentation](https://docs.github.com/en/code-security/security-advisories/about-security-advisories)
- **PGP Encryption** (for sensitive reports):
  - **Public Key**: [GPG key](https://keys.openpgp.org/vks/v1/by-fingerprint/7EE22AF642772ADE86C4BDCBA74F2AC73621FA84)
  - **Key ID**: `7EE2 2AF6 4277 2ADE 86C4  BDCB A74F 2AC7 3621 FA84`
  - [Instructions for encrypting messages](https://example.com/pgp-instructions)

### Guidelines for Reporting

When reporting a vulnerability:

1. Provide a detailed description of the issue, including steps to reproduce,
affected versions, and potential impacts.
2. Share any supporting information, such as logs or screenshots, to assist
in the investigation.
3. Avoid testing on systems you do not own or have explicit permission to test.

We value and appreciate responsible vulnerability disclosures.

## Coordinated Vulnerability Disclosure Policy

To protect our users:

- **Do not disclose vulnerabilities publicly** until a fix or mitigation is available.
- Allow the nogocomments team time to investigate, develop, and deploy a fix.
- Work collaboratively with the nogocomments team during active exploitation
scenarios to provide timely public guidance if necessary.

## Incident Response

Upon receiving a security report, the nogocomments team will:

1. **Acknowledge receipt** of your report within **48 hours**.
2. **Investigate the issue** and assess its severity.
3. **Develop a fix or mitigation plan** in collaboration with you,
if applicable.
4. **Coordinate disclosure timelines** with you to ensure user safety.
5. **Provide updates** at least every **7 days** until the issue is resolved.

## Recognition and Transparency

We recognize the efforts of researchers who responsibly disclose
vulnerabilities:

- With your permission, we will credit you in the release notes or other
public acknowledgments.
- Recognition is optional and can be declined if preferred.

## Dependency Vulnerabilities

If a vulnerability is identified in a dependency:

- We will coordinate with the dependency's maintainers and encourage 
private reporting through appropriate channels.
- Our goal is to protect users while respecting the vulnerability
disclosure timelines of upstream maintainers.

## Notification and Public Disclosure

Once a vulnerability is resolved, we will:

- **Publish release notes** detailing the resolution.
- **Update this `SECURITY.md` file**, if applicable.
- **Notify users** through our official communication channels (e.g.,
GitHub releases, mailing lists, or our website).

## Contributing to Security

We welcome contributions aimed at improving the security of
nogocomments. If you have ideas or patches to enhance security, please
submit a pull request.

## Legal and Testing Permissions

By submitting a vulnerability report or patch, you agree to the terms of
the nogocomments [LICENSE](LICENSE).

We request that security testing be performed only on instances of
nogocomments that you own or have explicit permission to test.
Unauthorized testing on third-party systems is strictly prohibited.

## Commitment to Security

The nogocomments team is committed to ensuring the safety and
integrity of our software and the trust of our users. By following this
policy, you contribute to a safer software ecosystem. Thank you for your
collaboration and commitment to security.
