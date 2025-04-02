#!/usr/bin/env bash
set -euo pipefail

issuer_url=$(yq .kubernetes.oidc.issuer_url config.yaml)
client_id=$(yq .kubernetes.oidc.client_id config.yaml)

echo -e "Using issuer: \e[1;35m$issuer_url\e[0m"
kubelogin get-token \
  --oidc-issuer-url="$issuer_url" \
  --oidc-client-id="$client_id" \
  --oidc-extra-scope="email"
