# Quickstart

effx provides this command line tool that facilates validation and synchronization of effx resources in your development environment. This CLI is the foundational tool in our git-based integrations.

![GitHub](https://img.shields.io/github/license/effxhq/effx-cli)


### Create your effx.yaml configurations

We support the creation of services, users, and teams via `effx.yaml` files using the `effx-cli` tool.

- Examples for users, teams, and multiple objects can be found [here](./examples/).
- A config must end with `effx.yaml`. e.g. `*effx.yaml`

### Usage

```bash
go run effx.go lint -d .
go run effx.go sync -d . -k ${EFFX_API_KEY}
go run effx.go event --title "title" --message "meassage" --service "dooku" --tags "key:value" --actions "alert:name:https://pagerduty.com -k $EFFX_API_KEY 
```

### disable automatic language and version detection

```bash
go run effx.go sync -d . -k ${EFFX_API_KEY} --disable-languauge-detection
```

### configure automatic service detection

you can configure which directory names can 
automatically detected services. For example any directories inside `services/` or `apps/` will contain detected services

```bash
export INFERRED_SERVICE_DIRECTORY_NAMES="services,apps"
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
