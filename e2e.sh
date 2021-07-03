#!/bin/sh
set -eux

stackit up --stack-name whodunnit-e2e --template e2e.yml
