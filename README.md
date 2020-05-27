Hosted at:
[dockerhub](https://hub.docker.com/r/effx/effx-cli)

# Quickstart
### Create your effx.yaml configurations
We support the creation of services, users, and teams via `effx.yaml` files using the `effx-cli` tool.

* The most common example is a [service definiton](./examples/service_effx.yaml)
* More examples for users, teams, and multiple objects can be found [here](./examples/).
* A config must end with `effx.yaml`. e.g. `*effx.yaml`

### Github Actions:
We've created two Github Actions for making it easy to automatically lint your config and sync it to our platform:
* Lint action setup [instructions](https://github.com/effxhq/effx-lint-action) from the action's Github repo. 
* Sync action setup [instructions](https://github.com/effxhq/effx-sync-action) from the action's Github repo.

### Gitlab CI:
Setup:
1. [grab your effx api key](https://app.effx.com/account_settings)
2. add `EFFX_API_KEY` variable with your api key in your `gitlab repo > settings > ci_cd > variables` page
3. add the following jobs to your `.gitlab-ci.yml` file

```yaml
# .gitlab-ci.yml

# Run a linter for your effx.yaml files
effx-lint:
  image: docker:latest
  variables:
    DOCKER_DRIVER: overlay
  services:
  - docker:dind
  stage: test
  script:
  - docker run -e EFFX_API_KEY --mount type=bind,source="$(pwd)",target=/app effx/effx-cli sync -d /app --dry-run

# Sync effx.yaml files, master branch only
effx-sync:
  image: docker:latest
  variables:
    DOCKER_DRIVER: overlay
  services:
  - docker:dind
  stage: deploy
  script:
  - docker run -e EFFX_API_KEY --mount type=bind,source="$(pwd)",target=/app effx/effx-cli sync -d /app
  only:
  - master

# Create an event every time this repo is deployed
effx-publish-deploy-event:
  image: docker:latest
  variables:
    DOCKER_DRIVER: overlay
  services:
  - docker:dind
  stage: deploy
  script:
  - |
    docker run -e EFFX_API_KEY effx/effx-cli event create \
      --name="$CI_PROJECT_TITLE deployed" \
      --desc="$CI_PROJECT_TITLE was deployed by $GITLAB_USER_EMAIL" \
      --service=$CI_PROJECT_TITLE \
      --integration_name=gitlab \
      --integration_version=1
  only:
  - master
```
