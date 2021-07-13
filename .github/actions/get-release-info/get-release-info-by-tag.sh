#!/bin/bash
# Gets the release info from a give tag
# First arg is repo name, second arg is tag, third arg is an oAuth token
# Echos out information for GH actions to pick up

if [ -z $GITHUB_API_URL ]; then
    GITHUB_API_URL="https://api.github.com"
fi

RELEASE_INFO=$(curl -sS -H "Accept: application/vnd.github.v3+json" -H "Authorization: Bearer $3" "$GITHUB_API_URL/repos/$1/releases/tags/$2")

if [ "$DEBUG" != "" ] && [ "$DEBUG" != "0" ]; then
    printf '%s\n' "$RELEASE_INFO"
fi 

ID=$(printf '%s' "$RELEASE_INFO" | jq -c -r ".id")
TITLE=$(printf '%s' "$RELEASE_INFO" | jq -c -r ".name")
PRERELEASE=$(printf '%s' "$RELEASE_INFO" | jq -c -r ".prerelease")
BODY=$(printf '%s' "$RELEASE_INFO" | jq -c -r ".body")

printf '::set-output name=id::%s\n' "$ID"
printf '::set-output name=title::%s\n' "$TITLE"
printf '::set-output name=pre-release::%s\n' "$PRERELEASE"
printf '::set-output name=body::%s\n' "$BODY"
