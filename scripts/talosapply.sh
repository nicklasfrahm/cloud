#!/usr/bin/env bash

set -euo pipefail

cluster_dir=""
cluster_name=""

preflight_check() {
  if [ "$#" -ne 1 ]; then
    echo "ℹ️  Usage: $0 <cluster-directory>"
    exit 1
  fi

  if [ ! -d "$1" ]; then
      echo "❌ Error: failed to find cluster directory: $1"
      exit 1
  fi

  dependencies=(talosctl yq)
  for cmd in "${dependencies[@]}"; do
    if ! command -v "$cmd" &> /dev/null; then
      echo "❌ Error: failed to find required command: $cmd"
      exit 1
    fi
  done
}

check_required_files() {
  local required_files=("meta.yaml" "controlplane.yaml" "worker.yaml" "talosconfig")

  for file in "${required_files[@]}"; do
    if [ ! -f "$cluster_dir/$file" ]; then
      echo "❌ Error: failed to find required file: $cluster_dir/$file"
      echo "ℹ️  Hint: run talosgen.sh first to generate the configuration files"
      exit 1
    fi
  done
}

apply_config_to_node() {
  local host="$1"
  local name="$2"
  local config_file="$3"
  local talosconfig="$cluster_dir/talosconfig"

  echo "✅ Info: Applying configuration (${host}/${name:-unknown})"

  extra_args=()
  if [ -n "$name" ]; then
    extra_args+=("-p" "[{\"op\": \"replace\", \"path\": \"/machine/network/hostname\", \"value\": \"$name\"}]")
  fi

  # Use timeout to handle unreachable nodes.
  if timeout 10s talosctl apply \
    --talosconfig "$talosconfig" \
    --endpoints "$host" \
    --nodes "$host" \
    --file "$config_file" \
    "${extra_args[@]}"; then
    echo "✅ Info: Successfully applied configuration (${host}/${name:-unknown})"
    return 0
  fi

  # Try to apply the config in maintenance mode if the normal apply failed.
  echo "⚠️  Warn: Normal apply failed, trying maintenance mode (${host}/${name:-unknown})"
  if timeout 10s talosctl apply \
    --talosconfig "$talosconfig" \
    --endpoints "$host" \
    --nodes "$host" \
    --file "$config_file" \
    --insecure \
    "${extra_args[@]}"; then
    echo "✅ Info: Successfully applied configuration in maintenance mode (${host}/${name:-unknown})"
    return 0
  fi

  echo "❌ Error: Failed to apply configuration (${host}/${name:-unknown})"
  return 1
}

apply_talos_configs() {
  local failed_nodes=0

  cluster_name=$(yq eval '.name' "$cluster_dir/meta.yaml" || echo "INVALID")
  if [ "$cluster_name" == "INVALID" ]; then
    echo "❌ Error: failed to read cluster name from meta.yaml"
    exit 1
  fi

  echo "Applying Talos configurations for cluster: $cluster_name"
  echo "=================================================="

  # Get the number of control plane and worker nodes.
  controlplane_count=$(yq '.nodes.controlplanes | length' "$cluster_dir/meta.yaml" 2>/dev/null || echo "0")
  for ((i=0; i<controlplane_count; i++)); do
    host=$(yq -r ".nodes.controlplanes[$i].host" "$cluster_dir/meta.yaml" 2>/dev/null || echo "")
    if [ -z "$host" ]; then
      echo "⚠️  Warn: Skipping node: failed to find host for control plane node: $i"
      failed_nodes=$((failed_nodes + 1))
      continue
    fi

    name=$(yq -r ".nodes.controlplanes[$i].name" "$cluster_dir/meta.yaml" 2>/dev/null || echo "")
    if [ -z "$name" ]; then
      echo "⚠️  Warn: Using auto-generated name: failed to find name for control plane node: $i"
    fi

    if ! apply_config_to_node "$host" "$name" "$cluster_dir/controlplane.yaml"; then
      failed_nodes=$((failed_nodes + 1))
    fi
  done

  worker_count=$(yq '.nodes.workers | length' "$cluster_dir/meta.yaml" 2>/dev/null || echo "0")

  for ((i=0; i<worker_count; i++)); do
    host=$(yq -r ".nodes.workers[$i].host" "$cluster_dir/meta.yaml" 2>/dev/null || echo "")
    if [ -z "$host" ]; then
      echo "⚠️  Warn: Skipping node: failed to find host for worker node: $i"
      failed_nodes=$((failed_nodes + 1))
      continue
    fi

    name=$(yq -r ".nodes.workers[$i].name" "$cluster_dir/meta.yaml" 2>/dev/null || echo "")
    if [ -z "$name" ]; then
      echo "⚠️  Warn: Using auto-generated name: failed to find name for worker node: $i"
    fi

    if ! apply_config_to_node "$host" "$name" "$cluster_dir/worker.yaml"; then
      failed_nodes=$((failed_nodes + 1))
    fi
  done

  total_nodes=$((controlplane_count + worker_count))
  echo "=================================================="
  echo "Total nodes:       $total_nodes"
  echo "Failed nodes:      $failed_nodes"
  echo "Successful nodes:  $((total_nodes - failed_nodes))"

  if [ $failed_nodes -gt 0 ]; then
    echo ""
    echo "⚠️  Some nodes failed to receive configuration updates."
    echo "   This is likely due to nodes being offline or unreachable."
    echo "   You can retry this script later for the failed nodes."
    exit 1
  else
    echo ""
    echo "✅ All nodes successfully configured!"
  fi
}

main() {
  preflight_check "$@"

  # Strip trailing slash if any.
  cluster_dir="${1%/}"

  check_required_files

  # Set the cluster name.
  cluster_name=$(yq '.name' "$cluster_dir/meta.yaml")

  apply_talos_configs
}

main "$@"
