curl -X PUT "https://api.effx.io/v2/services" -H "accept: application/json" -H "X-Effx-Api-Key: $EFFX_API_KEY" -H "Content-Type: application/json" -d'
{
  "version": "effx/v1",
  "kind": "service",
  "spec": {
    "name": "example-service2",
    "description": "this service is used for testing stuff",
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
      },
      "onCall": {
        "url": "https://www.pagerduty.com/platform/on-call-management/",
        "label": "pagerduty"
      }
    },
    "linkGroups": [
      {
        "label": "dashboards",
        "links": [
          {
            "url": "https://datadog.com",
            "label": "datadog"
          },
          {
            "url": "https://newrelic.com",
            "label": "newrelic"
          }
        ]
      },
      {
        "label": "runbook",
        "links": [
          {
            "url": "https://notion.so",
            "label": "deploy guide"
          }
        ]
      },
      {
        "label": "version control",
        "links": [
          {
            "url": "https://github.com/effxhq/.github",
            "label": "github"
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