#!/usr/bin/env bash

set -euo pipefail

preflight_check() {
  if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <cluster-directory>"
    exit 1
  fi

  if [ ! -d "$1" ]; then
      echo "Error: failed to find cluster directory: $1"
      exit 1
  fi

  dependencies=(talosctl yq find)
  for cmd in "${dependencies[@]}"; do
    if ! command -v "$cmd" &> /dev/null; then
      echo "Error: failed to find required command: $cmd"
      exit 1
    fi
  done
}

ensure_secret_bundle() {
  local cluster_dir="$1"
  local secret_bundle="$cluster_dir/secrets.yaml"

  if [ ! -f "$secret_bundle" ]; then
    talosctl gen secrets --output-file "$secret_bundle"
  fi
}

generate_talos_configs() {
  local cluster_dir="$1"

  # Find all control plane patches.
  mapfile -t control_plane_patch_files < <(find "$cluster_dir/patches/cp" -type f -name '*.yaml')
  # Find all worker patches.
  mapfile -t worker_patch_files < <(find "$cluster_dir/patches/worker" -type f -name '*.yaml')
  # Find all global patches.
  mapfile -t global_patch_files < <(find "$cluster_dir/patches/all" -type f -name '*.yaml')

  # Prepend an "@" to each patch file and separate with commas.
  control_plane_patch_flags=()
  for i in "${!control_plane_patch_files[@]}"; do
    control_plane_patch_flags+=("--config-patch-control-plane" "@${control_plane_patch_files[$i]}")
  done

  worker_patch_flags=()
  for i in "${!worker_patch_files[@]}"; do
    worker_patch_flags+=("--config-patch-worker" "@${worker_patch_files[$i]}")
  done

  global_patch_flags=()
  for i in "${!global_patch_files[@]}"; do
    global_patch_flags+=("--config-patch" "@${global_patch_files[$i]}")
  done

  cluster_name=$(yq e '.name' "$cluster_dir/meta.yaml" || echo "INVALID")
  if [ "$cluster_name" == "INVALID" ]; then
    echo "Error: failed to read cluster name: failed to read config: meta.yaml"
    exit 1
  fi

  cluster_endpoint=$(yq e '.endpoint' "$cluster_dir/meta.yaml" || echo "INVALID")
  if [ "$cluster_endpoint" == "INVALID" ]; then
    echo "Error: failed to read cluster endpoint: failed to read config: meta.yaml"
    exit 1
  fi

  # Generate Talos config for control plane nodes.
  talosctl gen config "$cluster_name" "$cluster_endpoint" \
    --with-secrets "$cluster_dir/secrets.yaml" \
    --output "$cluster_dir" \
    "${control_plane_patch_flags[@]}" \
    "${worker_patch_flags[@]}" \
    "${global_patch_flags[@]}" \
    --force
}

main() {
  preflight_check "$@"

  # Strip trailing slash if any.
  local cluster_dir="${1%/}"

  ensure_secret_bundle "$cluster_dir"
  generate_talos_configs "$cluster_dir"
}

main "$@"
