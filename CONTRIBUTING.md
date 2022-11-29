# Contributing

## Releases

### Preparation

Ensure that:
1. All PRs to be included in the release have been merged.
1. `CHANGELOG.md` details all changes relevant to end users and that PR links are correct.
1. The release date in `CHANGELOG.md` is correct.
1. The last merge to `master` (or relevant major version branch) built and ran all tests successfully.

#### Performing the release

1. On GitHub, 'Draft a new release':
    1. Tag version - of the form v1.2.3
    1. Target - generally `master` unless the release is a minor/patch for a previous major version for which we have a branch.
    1. Release title - as the Tag version
    1. Description - copy directly from `CHANGLEOG.md`, ensuring that the formatting looks correct in the preview.
1. Publish release
1. Update and push to `next`