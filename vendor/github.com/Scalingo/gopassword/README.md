# Go Password v1.1.0

Simple password generator in Go. Use `crypto/rand`

```go
// Password of 64 characters
gopassword.Generate()

// Password of 42 characters
gopassword.Generate(42)
```

## Release a New Version

Bump new version number in `CHANGELOG.md` and `README.md`.

Commit, tag and create a new release:

```sh
version="1.1.0"
git switch --create release/${version}
git add CHANGELOG.md README.md
git commit -m "Bump v${version}"
git push --set-upstream origin release/${version}
gh pr create --reviewer=EtienneM --title "$(git log -1 --pretty=%B)"
```

Once the pull request merged, you can tag the new release.

```sh
git tag v${version}
git push origin master v${version}
gh release create v${version}
```

The title of the release should be the version number and the text of the release is the same as the changelog.
