#!/usr/bin/env bash

# Exit on error
set -eou pipefail

# Check for required arguments
if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <GitHub App ID> <Path to PEM file>"
  exit 1
fi

APP_ID="$1"
PEM_FILE="$2"

# Generate a JWT for the GitHub App using OpenSSL
HEADER=$(echo -n '{"alg":"RS256","typ":"JWT"}' | base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')
PAYLOAD=$(echo -n "{\"iat\":$(date +%s),\"exp\":$(($(date +%s) + 600)),\"iss\":\"$APP_ID\"}" | base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')
SIGNATURE=$(echo -n "$HEADER.$PAYLOAD" | openssl dgst -sha256 -sign "$PEM_FILE" | base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')
JWT="$HEADER.$PAYLOAD.$SIGNATURE"

# Fetch installation IDs
curl -s -H "Authorization: Bearer $JWT" -H "Accept: application/vnd.github+json" \
  https://api.github.com/app/installations | jq '.[].id'
