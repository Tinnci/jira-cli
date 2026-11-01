# Tinnci fork

This fork of [ankitpokhrel/jira-cli](https://github.com/ankitpokhrel/jira-cli) adds what
SPM's project setup needs from Jira Server 9.17. SPM reads Jira only through this CLI,
which keeps the token.

- `jira release list --raw`: versions as JSON. Jira Server and Data Center
  (`installation: Local`) now use `/rest/api/2/project/{key}/versions`; upstream calls the
  Cloud-only v3 endpoint and gets an HTML 404.
- `jira filter list [--raw]`: the user's starred filters with their JQL
  (`/rest/api/2/filter/favourite`).
- `jira issue query --jql Q [--fields a,b] [--from N] [--limit N]`: the v2 `/search`
  response unchanged, with `total` and custom fields. `issue list --raw` drops fields that
  jira-cli's own types do not model, such as a severity custom field.
- `jira issue count --jql Q [--jql Q2 …]`: `[{"jql","total"}]` without reading issues.

Measured against PLV on 2026-11-01: 47 versions in 0.09 s, two counts in 0.1 s, and 1000
issues with four fields in one 1 s call (upstream needs one `issue view` per issue).
