#!/bin/sh
set -eux

stackit up --stack-name aggregation-bus --template aggregation-bus.yml
#stackit up --stack-name whodunnit-rules --template org-wide-rules.yml
