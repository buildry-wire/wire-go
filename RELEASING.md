# Releasing wire-go

Go has no package registry — a release is a semver git tag. pkg.go.dev indexes the
module the first time the module proxy fetches the version.

## Cut a release
1. Update `CHANGELOG.md`: move items under a new `## [x.y.z] - YYYY-MM-DD` section.
2. Commit on `main`.
3. Tag and push:
   ```bash
   git tag vX.Y.Z
   git push origin vX.Y.Z
   ```
4. The `release` workflow runs tests, warms `proxy.golang.org`, and creates a GitHub Release.

## Verify
```bash
go get github.com/buildry-wire/wire-go@vX.Y.Z
```
Then check https://pkg.go.dev/github.com/buildry-wire/wire-go@vX.Y.Z (may take a few minutes).

No secrets or registry accounts are required for Go.
