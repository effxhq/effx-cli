Hosted at:
[dockerhub](https://hub.docker.com/r/effx/effx-cli)

# Quickstart
### effx.yaml
* it must be named `effx.yaml`
* find an example [here](./examples/effx.yaml)

### Github Action:
Follow the setup [instructions](https://github.com/effxhq/effx-sync-action) from the action's github repo.

### Gitlab CI:
setup:
1. [grab your effx api key](https://app.effx.com/account_settings)
2. add `EFFX_API_KEY` variable with your api key in your `gitlab repo > settings > ci_cd > variables` page
3. add the following jobs to your `.gitlab.ci.yml` file

```yaml
# .gitlab.ci.yml

# run a linter for your effx.yaml files
effx-lint:
  image: docker:latest
  variables:
    DOCKER_DRIVER: overlay
  services:
  - docker:dind
  stage: test
  script:
  - docker run -e EFFX_API_KEY --mount type=bind,source="$(pwd)",target=/app effx/effx-cli sync -d /app --dry-run

# sync effx.yaml files, master branch only
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

# create an event every time this repo is deployed
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
