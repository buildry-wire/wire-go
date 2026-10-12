# Changelog

Format: Keep a Changelog; semver.

## [Unreleased]

## [1.0.1] - 2026-10-12

### Changed
- README: amounts are whole tögrög (MNT integer, no fractional part — `50000` is 50,000 ₮, not 500.00). Examples no longer use the `sandbox` operator or test mode, which the production API does not offer; omitting `allowed_operators` routes to the operators enabled on your project.

## [1.0.0] - 2026-06-07
First public release.

### Added
- Client with auth, idempotency, retry/backoff.
- PaymentIntents, Charges, Events, WebhookEndpoints services.
- Auto-pagination iterator.
- Webhook signature verification.
