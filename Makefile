#################
### Cloud CLI ###
#################

SOURCES		:= $(shell find . -type f -name '*.go')
VERSION		:= $(shell git describe --tags --always --dirty)
GO_FLAGS	:= -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: install
install: bin/cloud
	sudo install -Dm755 bin/cloud /usr/local/bin/cloud

bin/cloud: $(SOURCES)
	CGO_ENABLED=0 go build $(GO_FLAGS) -o bin/cloud ./cmd/cloud/main.go

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
