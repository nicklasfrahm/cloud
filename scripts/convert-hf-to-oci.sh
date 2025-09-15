#!/usr/bin/env bash
set -euo pipefail

# Script: convert-hf-to-oci.sh
# Purpose: Convert Hugging Face models to OCI images for use with KServe
#          Automatically detects and converts .nemo files to PyTorch format
# Usage: ./convert-hf-to-oci.sh <huggingface-repo> <oci-repo> [model-name]
#
# Arguments:
#   huggingface-repo  - Hugging Face repository name (e.g., gpt2, bert-base-uncased, nvidia/parakeet-tdt-0.6b-v3)
#   oci-repo          - OCI repository with tag (e.g., ghcr.io/myorg/gpt2-hf:latest)
#   model-name        - Optional model directory name (default: model)
#
# Examples:
#   ./convert-hf-to-oci.sh "gpt2" "ghcr.io/myorg/gpt2-hf:latest"
#   ./convert-hf-to-oci.sh "nvidia/parakeet-tdt-0.6b-v3" "ghcr.io/myorg/parakeet:v3" "parakeet"

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}info:${NC} $*"
}

log_success() {
    echo -e "${GREEN}success:${NC} $*"
}

log_warning() {
    echo -e "${YELLOW}warning:${NC} $*"
}

log_error() {
    echo -e "${RED}error:${NC} $*" >&2
}

# Usage function
show_usage() {
    cat << EOF
Usage: $(basename "$0") <huggingface-repo> <oci-repo> [model-name]

Arguments:
  huggingface-repo  Hugging Face repository name (e.g., gpt2, nvidia/parakeet-tdt-0.6b-v3)
  oci-repo          OCI repository with tag (e.g., ghcr.io/myorg/gpt2-hf:latest)
  model-name        Optional model directory name (default: model)

Examples:
  $(basename "$0") "gpt2" "ghcr.io/myorg/gpt2-hf:latest"
  $(basename "$0") "nvidia/parakeet-tdt-0.6b-v3" "ghcr.io/myorg/parakeet:v3" "parakeet"

This script downloads a Hugging Face model and packages it as an OCI image
for use with the HuggingFace runtime in KServe. If .nemo files are detected,
they will be automatically converted to PyTorch format for broader compatibility.
EOF
}

# Check for help flags
if [[ $# -eq 0 ]] || [[ "${1:-}" == "-h" ]] || [[ "${1:-}" == "--help" ]]; then
    show_usage
    exit 0
fi

# Validate arguments
if [[ $# -lt 2 ]]; then
    log_error "Missing required arguments"
    show_usage
    exit 1
fi

HF_REPO="${1}"
OCI_REPO="${2}"
MODEL_NAME="${3:-model}"

# Create venv in current working directory
VENV_DIR="${PWD}/.venv"

# Create temporary directory
TMP_DIR="$(mktemp -d)"
MODEL_DIR="${TMP_DIR}/models/${MODEL_NAME}"

# Cleanup function
cleanup() {
    local exit_code=$?
    if [[ -n "${TMP_DIR:-}" ]] && [[ -d "${TMP_DIR}" ]]; then
        log_info "Cleaning up temporary directory: ${TMP_DIR}"
        rm -rf "${TMP_DIR}"
    fi
    if [[ $exit_code -ne 0 ]]; then
        log_error "Script failed with exit code $exit_code"
    fi
    exit $exit_code
}

# Set trap for cleanup
trap cleanup EXIT INT TERM

log_info "Downloading Hugging Face repository: ${HF_REPO}"
git clone "https://huggingface.co/${HF_REPO}" "${MODEL_DIR}"

# Check for .nemo files and convert to PyTorch if found
if find "${MODEL_DIR}" -name "*.nemo" -type f | grep -q .; then
    log_info "Found .nemo file(s), converting to PyTorch format"

    # Create or reuse persistent virtual environment for NeMo conversion
    if [[ ! -d "${VENV_DIR}" ]]; then
        log_info "Creating persistent virtual environment: ${VENV_DIR}"
        python3 -m venv "${VENV_DIR}"

        # Install NeMo toolkit using the venv's pip directly
        log_info "Installing NeMo Toolkit in virtual environment"
        "${VENV_DIR}/bin/pip3" install --upgrade pip
        "${VENV_DIR}/bin/pip3" install nemo_toolkit[all]
    else
        log_info "Reusing existing virtual environment: ${VENV_DIR}"
        # Check if NeMo is installed, install if missing
        if ! "${VENV_DIR}/bin/python3" -c "import nemo" 2>/dev/null; then
            log_info "NeMo not found in existing venv, installing..."
            "${VENV_DIR}/bin/pip3" install --upgrade pip
            "${VENV_DIR}/bin/pip3" install nemo_toolkit[all]
        fi
    fi

    # Find all .nemo files and convert them
    find "${MODEL_DIR}" -name "*.nemo" -type f | while read -r nemo_file; do
        log_info "Converting: $(basename "${nemo_file}")"

        # Run conversion using the venv's python directly
        "${VENV_DIR}/bin/python3" << EOF
import nemo
import os

nemo_file = "${nemo_file}"
model_dir = "${MODEL_DIR}"

# Load the NeMo model
model = nemo.collections.nlp.models.language_modeling.megatron_gpt_model.MegatronGPTModel.restore_from(nemo_file)

# Save as PyTorch model
pytorch_dir = os.path.join(model_dir, "pytorch_model")
os.makedirs(pytorch_dir, exist_ok=True)

# Save the model state dict and config
import torch
torch.save(model.state_dict(), os.path.join(pytorch_dir, "pytorch_model.bin"))

# Save tokenizer if available
if hasattr(model, 'tokenizer') and model.tokenizer is not None:
    model.tokenizer.save_pretrained(pytorch_dir)

print(f"Converted {nemo_file} to PyTorch format in {pytorch_dir}")
EOF
    done

    # Clean up original .nemo files and keep only converted PyTorch files
    log_info "Cleaning up original .nemo files"
    find "${MODEL_DIR}" -name "*.nemo" -type f -delete

    # If pytorch_model directory exists, move its contents to the main model directory
    if [[ -d "${MODEL_DIR}/pytorch_model" ]]; then
        log_info "Moving converted PyTorch files to model root"
        mv "${MODEL_DIR}/pytorch_model"/* "${MODEL_DIR}/"
        rmdir "${MODEL_DIR}/pytorch_model"
    fi

    log_success "NeMo model(s) converted to PyTorch format"
else
    log_info "No .nemo files found, proceeding with standard model packaging"
fi

log_info "Creating Dockerfile (model-only image)"
cat > "${TMP_DIR}/Dockerfile" <<EOF
FROM alpine:latest
COPY models /models
EOF

log_info "Building OCI image: ${OCI_REPO}"
docker build -t "${OCI_REPO}" -f "${TMP_DIR}/Dockerfile" "${TMP_DIR}"

log_info "Pushing to OCI registry: ${OCI_REPO}"
docker push "${OCI_REPO}"

log_success "Hugging Face model published as OCI image: ${OCI_REPO}"
log_success "Use this image with the HuggingFace runtime in KServe."
