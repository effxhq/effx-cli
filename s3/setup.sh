curl -i -X PUT "https://api.effx.io/v2/services" -H "accept: application/json" -H "X-Effx-Api-Key: $EFFX_API_KEY" -H "Content-Type: application/json" -d'
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
            "url": "http://effx.com/help",
            "label": "Getting Started"
          },
          {
            "url": "https://effxhq.github.io/effx-api-v2/",
            "label": "API"
          },
          {
            "url": "https://effx.com/help/?shell#cli-documentation",
            "label": "effx CLI"
          },
          {
            "url": "https://github.com/effxhq/vcs-connect",
            "label": "vcs-connect"
          }
        ]
      },
      {
        "label": "dashboards",
        "links": [
          {
            "url": "http://status.effx.com",
            "label": "status"
          },
          {
            "url": "https://p.datadoghq.com/sb/f655e3381-6bdce95d9ec6f133df05532b21b5f5dc",
            "label": "datadog"
          }
        ]
      },
      {
        "label": "runbook",
        "links": [
          {
            "url": "https://effx.com/blog/how-to-write-a-runbook",
            "label": "best practices"
          }
        ]
      },
      {
        "label": "version control",
        "links": [
          {
            "url": "https://github.com/effxhq/.github",
            "label": "github workflows"
          },
          {
            "url": "https://github.com/effxhq/vcs-connect",
            "label": "vcs-connect"
          }
        ]
      },
      {
        "label": "continuous integration",
        "links": [
          {
            "url": "https://effx.com/help/?yaml#using-github-actions",
            "label": "github actions"
          },
          {
            "url": "https://effx.com/help/?yaml#using-gitlab-ci-cd-jobs",
            "label": "gitlab"
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
