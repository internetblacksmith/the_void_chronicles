# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of The Void Chronicles SSH Reader seriously. If you have discovered a security vulnerability, please follow these steps:

### 1. Do NOT Create a Public Issue

Security vulnerabilities should **never** be reported via public GitHub issues as this could put users at risk.

### 2. Email Us Privately

Send details of the vulnerability to: **security@voidchronicles.space**

Please include:
- Type of vulnerability
- Full paths of source file(s) related to the issue
- Location of the affected source code (tag/branch/commit or direct URL)
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

### 3. Wait for Initial Response

You should receive an initial response within 48 hours acknowledging receipt of your report.

### 4. Disclosure Timeline

- **Day 0**: You report the vulnerability
- **Day 1-2**: We acknowledge receipt
- **Day 3-7**: We investigate and validate the issue
- **Day 7-30**: We develop and test a fix
- **Day 30-45**: We release the fix
- **Day 45+**: Public disclosure (coordinated with reporter)

## Security Considerations

### SSH Server Security

The SSH reader uses password authentication by default. For production deployments:

1. **Use strong passwords**: Change the default password
2. **Restrict access**: Use firewall rules to limit access
3. **Monitor logs**: Regularly check access logs for suspicious activity
4. **Keep updated**: Apply security patches promptly

### Data Privacy

- **No telemetry**: The application does not collect or transmit user data
- **Local storage**: Reading progress is stored locally only
- **No external connections**: Besides SSH, no external network connections are made

### Known Security Features

- Password authentication for SSH access
- User isolation (each SSH user has separate progress tracking)
- Read-only book content (users cannot modify the source material)
- No shell access (SSH is restricted to the TUI application)

## Security Best Practices for Deployment

### Railway Deployment
- Set strong `SSH_PASSWORD` environment variable
- Use Railway's built-in DDoS protection
- Monitor access logs via Railway dashboard

### VPS Deployment
```bash
# Change default password
export SSH_PASSWORD="your-strong-password-here"

# Restrict SSH access by IP (iptables example)
iptables -A INPUT -p tcp --dport 23234 -s YOUR_IP/32 -j ACCEPT
iptables -A INPUT -p tcp --dport 23234 -j DROP

# Use fail2ban to prevent brute force attacks
apt-get install fail2ban
```

### Docker Deployment
```yaml
# docker-compose.yml security additions
services:
  void-reader:
    read_only: true
    security_opt:
      - no-new-privileges:true
    cap_drop:
      - ALL
    cap_add:
      - SETUID
      - SETGID
```

## Security Updates

Security updates will be released as:
- **Patch versions** (x.x.1) for non-breaking security fixes
- **Minor versions** (x.1.0) for security fixes that may break compatibility
- **Security advisories** via GitHub Security Advisories

## Acknowledgments

We appreciate responsible disclosure of security vulnerabilities. Security researchers who follow this policy will be acknowledged in our release notes and Hall of Fame (unless they prefer to remain anonymous).

## Contact

- Security issues: security@voidchronicles.space
- General inquiries: Open a GitHub issue
- Urgent issues: Include [URGENT] in email subject

Thank you for helping keep The Void Chronicles secure!