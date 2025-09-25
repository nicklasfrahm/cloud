#!/usr/bin/env bash
set -euo pipefail

issuer_url=$(yq -r .kubernetes.oidc.issuer_url config.yaml)
client_id=$(yq -r .kubernetes.oidc.client_id config.yaml)

echo -e "Using issuer: \e[1;35m$issuer_url\e[0m"
kubectl oidc-login get-token \
  --oidc-issuer-url="$issuer_url" \
  --oidc-client-id="$client_id" \
  --oidc-extra-scope="email" \
  --oidc-extra-scope="groups" | jq .
