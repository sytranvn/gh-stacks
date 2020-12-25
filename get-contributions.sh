#!/usr/bin/env bash

RUNNER=
JQ=
if command -v gh &> /dev/null
then
  RUNNER=gh
elif command -v hub &> /dev/null
then
  RUNNER=hub
fi

if command -v jq &> /dev/null
then
  JQ=jq
fi

if [[ -e "$RUNNER" ]]; then
  echo 'Depend on gh or hub. Please install github.com/cli/cli or github.com/github/hub'
  exit 1
fi

if [[ -e "$JQ" ]]; then
  echo 'Depend on jq. Please install jq'
  exit 1
fi

$RUNNER api graphql -F query=@contributions.gql | jq '
  .data.viewer.repositoriesContributedTo.nodes as $cons 
  | .data.viewer.repositories.nodes as $repos 
  | [$cons[], $repos[]] 
  | map(select(.languages.nodes | length  > 0))' > contributions.json

