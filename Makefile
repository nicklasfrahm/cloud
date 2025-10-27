################
### OpenTofu ###
################

OPENTOFU_ROOT_MODULE	?= deploy/opentofu
OPENTOFU_PLAN					?= $(OPENTOFU_ROOT_MODULE)/opentofu.tfplan
OPENTOFU_ROOT_SOURCES	?= $(shell find $(OPENTOFU_ROOT_MODULE) -maxdepth 1 -type f -name '*.tf')

$(OPENTOFU_PLAN): $(OPENTOFU_ROOT_SOURCES)
	tofu -chdir=deploy/opentofu init
	tofu -chdir=deploy/opentofu plan -out=opentofu.tfplan | tee tofu.log
	@sed -i 's/\x1b\[[0-9;]*m//g' tofu.log

opentofu-plan: $(OPENTOFU_PLAN) ## Plan the infrastructure changes.

.PHONY: opentofu-count
opentofu-count: opentofu-plan ## Count the number of changes in the plan.
	@tofu -chdir=deploy/opentofu show -json opentofu.tfplan | jq -r '.resource_changes[].change.actions | join(",")' | grep -Ecv '^no-op$$' || true

.PHONY: opentofu-apply
opentofu-apply: ## Apply the infrastructure changes.
	tofu -chdir=deploy/opentofu init
	tofu -chdir=deploy/opentofu apply -auto-approve

#####################
### LLM Packaging ###
#####################

REGISTRY							?= ghcr.io
NAMESPACE							?= $(shell whoami)
MODEL_REPO						?= "Qwen/Qwen2.5-0.5B-Instruct"
MODEL_ALIAS						?= "qwen2-instruct"

.PHONY: llm
llm: ## Generate LLM assets.
	docker build \
		--secret id=HF_TOKEN \
		-f llm.Containerfile \
		--build-arg MODEL="$(MODEL_REPO)" \
		--build-arg ALIAS="$(MODEL_ALIAS)" \
		-t $(REGISTRY)/$(NAMESPACE)/models/$(MODEL_ALIAS):latest \
		--load \
		-t $(MODEL_ALIAS):latest .
