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
3. add the following jobs to your `.gitlab.ci.yml` file

```yaml
# .gitlab.ci.yml

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

### Spinnaker:
Setup a Custom Webhook Stage [(ref)](https://spinnaker.io/guides/operator/custom-webhook-stages/):
1. [Grab your effx api key](https://app.effx.com/account_settings)
2. Replace your `EFFX_API_KEY` in the yaml snippet below and save it into your `orca-local.yaml` file. Replace the `DECK_HOSTNAME` with your Spinnaker URL.
3. Deploy the configuration update.
4. Use the new custom stage in the Spinnaker UI. Use the corresponding service name to link the event in Effx.

```yaml
---
webhook:
  preconfigured:
    - label: Effx - Deployment Event
      type: effxEvent
      enabled: true
      description: Post a Deployment Event to Effx
      method: PUT
      customHeaders:
        X-Effx-Api-Key:
          - <EFFX_API_KEY>
        Content-Type:
          - application/json
      url: https://api.effx.io/v1/events
      payload:
        '{
          "name": "Spinnaker - Webhook Stage",
          "description": "${execution.getName()} triggered by ${trigger.user} (${trigger.type})",
          "produced_at_time_milliseconds": ${execution.getBuildTime()},
          "actions": [
            {
              "level": "info",
              "name": "Pipeline",
              "url": "https://<DECK_HOSTNAME>/#/applications/${execution.getApplication()}/executions/${execution.getId()}"
            }
          ],
          "integration": {
            "name": "spinnaker",
            "image_url": "https://effx-post-services-integration-images.s3.us-west-2.amazonaws.com/spinnaker-logo.png"
          },
          "service": {
            "name": "${parameterValues["service"]}"
          }
        }'
      parameters:
        - label: Effx Service
          name: service
          description: The Effx Service for this Event
```
