# story: e19s03
# Title: Configurable Release Name and Body Templates

## User Story
As a repository owner,
I want to customize release title and body formats via Go templates,
So that release headings and footers match project branding.

## Acceptance Criteria
1. Support custom templates in `buildReleasePayload` with context variables.
2. Default template renders standard `v{version}` format.
