# Contributing

## Releases

### Preparation

Ensure that:
1. All PRs to be included in the release have been merged.
2. `CHANGELOG.md` details all changes relevant to end users and that PR links are correct.
3. The release date in `CHANGELOG.md` is correct.
4. The last merge to `main` (or relevant major version branch) built and ran all tests successfully.

#### Performing the release

1. Run `make VERSION=x.xx.x bump` to set the desired version number and date of release.
2. On GitHub, 'Draft a new release':
    1. Tag version - of the form `vx.xx.x` - This can be obtained from the `CHANGELOG.md`
    2. Target - generally `main` unless the release is a minor/patch for a previous major version for which we have a branch.
    3. Release title - as the Tag version
    4. Description - copy directly from `CHANGLEOG.md`, ensuring that the formatting looks correct in the preview.
3. Publish release
4. Update and push to `next`
5. Run `npm publish` to publish the NPM package.
