# Quickstart

### Create your effx.yaml configurations

We support the creation of services, users, and teams via `effx.yaml` files using the `effx-cli` tool.

- Examples for users, teams, and multiple objects can be found [here](./examples/).
- A config must end with `effx.yaml`. e.g. `*effx.yaml`

### Usage

```bash
go run effx.go lint -d .
go run effx.go sync -d . -k ${EFFX_API_KEY}
```

### Github Actions

We've created two Github Actions for making it easy to automatically lint your config and sync it to our platform:

- Lint action setup [instructions](https://github.com/effxhq/effx-lint-action) from the action's Github repo.
- Sync action setup [instructions](https://github.com/effxhq/effx-sync-action) from the action's Github repo.

### Gitlab

```yaml
image: ubuntu:latest

before_script:
  - curl -Lo effx https://effx.run/effx-cli/releases/latest/effx-cli_Linux_x86_64
  - sudo install effx /usr/local/bin

lint-all:
  stage: test
  script:
    - effx lint -d .

sync-all:
  stage: deploy
  variables:
    EFFX_API_KEY: ${EFFX_API_KEY}
  script:
    - effx sync -d .
  only:
    - master
```
