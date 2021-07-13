#!/bin/bash
# Script args are: [reponame] [tag] [path] [token]
if [ -z $GITHUB_API_URL ]; then
    GITHUB_API_URL="https://api.github.com"
fi

RELEASE_INFO=$(curl -sS -H "Accept: application/vnd.github.v3+json" -H "Authorization: Bearer $4" "$GITHUB_API_URL/repos/$1/releases/tags/$2")
UPLOAD_URL=$(printf '%s' "$RELEASE_INFO" | jq -c -r ".upload_url")
UPLOAD_URL=$(printf '%s' "$UPLOAD_URL" | sed 's/{?.*}$//')
ASSETS=$(printf '%s' "$RELEASE_INFO" | jq -c -r ".assets[] | {name, id}")

if [ "$DEBUG" != "" ] && [ "$DEBUG" != "0" ]; then
    printf 'UPLOAD URL: %s\n' "$UPLOAD_URL"
    printf 'ASSET NAMES:\n%s\n' "$ASSETS"
fi

# For each asset in path, we want to:
# Delete it's current version if it exists (asset with the same name of the file is already there)
# Upload the asset using the UPLOAD_URL
for f in "$3"/*
do
    BASENAME_EXT=${f##*/}
    
    printf 'processing %s (%s)\n' "$f" "$BASENAME_EXT"
    
    MATCHING_ASSET_ID=$(printf '%s' "$ASSETS" | jq -c -r --arg x "$BASENAME_EXT" '. | select(.name == $x) | .id')
    if [ "$MATCHING_ASSET_ID" != "" ]; then
        printf 'Deleting asset %s (id %s)\n' "$BASENAME_EXT" "$MATCHING_ASSET_ID"
        curl -sS -X DELETE -H "Accept: application/vnd.github.v3+json" -H "Authorization: Bearer $4" \
        "$GITHUB_API_URL/repos/$1/releases/assets/$MATCHING_ASSET_ID" | jq || exit 1
    fi

    printf 'Uploading asset %s\n' "$BASENAME_EXT"

    ENCODED_NAME=$(jq -c -r -n --arg x "$BASENAME_EXT" '$x|@uri')
    curl -sS -X POST -H "Accept: application/vnd.github.v3+json" -H "Authorization: Bearer $4" \
       -H "Content-Type: application/octet-stream" \
       --data-binary "@$f" "$UPLOAD_URL?name=$ENCODED_NAME" | jq || exit 1
done
