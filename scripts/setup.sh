curl -X PUT "https://api.effx.io/v2/services" -H "accept: application/json" -H "X-Effx-Api-Key: $EFFX_API_KEY" -H "Content-Type: application/json" -d'
{
  "version": "effx/v1",
  "kind": "service",
  "spec": {
    "name": "example-service",
    "description": "this service contains useful links for using Effx",
    "tags": {
      "go": "1.15.4",
      "group": "example"
    },
    "annotations": {
      "effx.io/owned-by": "example team"
    },
    "contact": {
      "email": "support@effx.com",
      "chat": {
        "label": "#community",
        "url": "https://join.slack.com/t/m11s/shared_invite/zt-j71p8afk-07MmPgrXlUd7qF_s6qXovg"
      }
    },
    "linkGroups": [
      {
        "label": "documentation",
        "links": [
          {
            "url": "https://help.effx.com",
            "label": "help"
          },
          {
            "url": "https://effxhq.github.io/effx-api-v2/",
            "label": "API"
          },
          {
            "url": "placeholder",
            "label": "vcs-connect"
          }
        ]
      },
      {
        "label": "dashboards",
        "links": [
          {
            "url": "https://status.effx.com",
            "label": "status"
          },
          {
            "url": "placeholder for wuher uptime & latency",
            "label": "datadog"
          }
        ]
      },
      {
        "label": "version control",
        "links": [
          {
            "url": "https://github.com/effxhq/.github",
            "label": "github workflows"
          }
        ]
      },
      {
        "label": "continuous integration",
        "links": [
          {
            "url": "https://github.com/effxhq/.github/blob/main/workflow-templates/effx-cli-ci.yml",
            "label": "github workflows"
          },
          {
            "url": "https://github.com/marketplace?type=actions&query=effxhq",
            "label": "github actions"
          },
          {
            "url": "https://circleci.com/developer/orbs/orb/effx/effx-cli",
            "label": "circleci Orb"
          }
        ]
      }
    ]
  }
}
'