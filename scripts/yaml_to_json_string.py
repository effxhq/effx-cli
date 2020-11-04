#!/usr/bin/env python3
"""
Print out a Yaml as JSON string for effunction data
"""

import json
import yaml


def print_config_body():
    """
    Print in config request format
    """
    with open("examples/team.effx.yml", "r") as fhandle:
        doc = fhandle.read()

    body = {
        "fileContents": doc,
        "tags": {
            "language": "go"
        },
        "annotations": {
            "effx.io/source": "effx-cli"
        }

    }
    print(json.dumps(json.dumps(body)))


def print_contents_only():
    """
    Print contents of yaml file as string escaped json
    """
    with open("examples/team.effx.yml", "r") as fhandle:
         doc = yaml.load(fhandle, Loader=yaml.Loader)

    print(json.dumps(json.dumps(doc)))

def main():
    """
    Convert an Effx Yaml to string
    """
    print_contents_only()


if __name__ == "__main__":
    main()
