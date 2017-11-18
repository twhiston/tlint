# tlint

My own dumb linting program, you probably shouldnt use this.
Its really inflexible, not configurable at all etc...

Runs:

- goimports
- gometalinter
- hadolint
- shellcheck
- checkmake

over everything under the cwd.
It will glob appropriately for the tools requirements