#!/usr/bin/env python3
"""
Print out a Yaml as JSON string for effunction data
"""

import yaml
import json


def main():
    """
    Convert an Effx Yaml to string
    """
    with open("examples/team.effx.yml", "r") as fhandle:
        doc = yaml.load(fhandle, Loader=yaml.Loader)
    print(json.dumps(json.dumps(doc)))

if __name__ == "__main__":
    main()
