BASE_PATH ?= $(CURDIR)
# Set to empty string to echo some command lines which are hidden by default.
SILENT ?= @

# GENERATED_API_XXX and PROTO_API_XXX variables contain standard paths used to
# generate gRPC proto messages, services for the API.
PROTO_BASE_PATH = $(CURDIR)/proto
ALL_PROTOS = $(shell find $(PROTO_BASE_PATH) -name '*.proto')
SERVICE_PROTOS = $(filter %_service.proto,$(ALL_PROTOS))

ALL_PROTOS_REL = $(ALL_PROTOS:$(PROTO_BASE_PATH)/%=%)


GENERATED_BASE_PATH = $(BASE_PATH)/generated
GENERATED_PB_SRCS = $(ALL_PROTOS_REL:%.proto=$(GENERATED_BASE_PATH)/%.pb.go)

##############
## Protobuf ##
##############
# Set some platform variables for protoc.
# If the proto version is changed, be sure it is also changed in qa-tests-backend/build.gradle.
PROTOC_VERSION := 3.20.1
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
PROTOC_ARCH = linux
endif
ifeq ($(UNAME_S),Darwin)
PROTOC_ARCH = osx
endif

PROTO_PRIVATE_DIR := $(BASE_PATH)/bin/protoc

PROTOC_DIR := $(PROTO_PRIVATE_DIR)/protoc-$(PROTOC_ARCH)-$(PROTOC_VERSION)

PROTOC := $(PROTOC_DIR)/bin/protoc

PROTOC_DOWNLOADS_DIR := $(PROTO_PRIVATE_DIR)/.downloads

PROTO_GOBIN := $(BASE_PATH)/bin

$(PROTOC_DOWNLOADS_DIR):
	@echo "+ $@"
	$(SILENT)mkdir -p "$@"

$(PROTO_GOBIN):
	@echo "+ $@"
	$(SILENT)mkdir -p "$@"

PROTOC_ZIP := protoc-$(PROTOC_VERSION)-$(PROTOC_ARCH)-x86_64.zip
PROTOC_FILE := $(PROTOC_DOWNLOADS_DIR)/$(PROTOC_ZIP)

$(PROTOC_FILE): $(PROTOC_DOWNLOADS_DIR)
	@echo "+ $@"
	$(SILENT)wget -q "https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC_ZIP)" -O "$@"

.PRECIOUS: $(PROTOC_FILE)

$(PROTOC):
	@echo "+ $@"
	$(SILENT)$(MAKE) "$(PROTOC_FILE)"
	$(SILENT)mkdir -p "$(PROTOC_DIR)"
	$(SILENT)unzip -q -o -d "$(PROTOC_DIR)" "$(PROTOC_FILE)"
	$(SILENT)test -x "$@"


PROTOC_INCLUDES := $(PROTOC_DIR)/include/google

PROTOC_GEN_GO_BIN := $(PROTO_GOBIN)/protoc-gen-gofast

MODFILE_DIR := $(PROTO_PRIVATE_DIR)/modules

$(MODFILE_DIR)/%/UPDATE_CHECK: go.sum
	@echo "+ Checking if $* is up-to-date"
	$(SILENT)mkdir -p $(dir $@)
	$(SILENT)go list -m -json $* | jq '.Dir' >"$@.tmp"
	$(SILENT)(cmp -s "$@.tmp" "$@" && rm "$@.tmp") || mv "$@.tmp" "$@"

$(PROTOC_GEN_GO_BIN): $(MODFILE_DIR)/github.com/gogo/protobuf/UPDATE_CHECK $(PROTO_GOBIN)
	@echo "+ $@"
	$(SILENT)GOBIN=$(PROTO_GOBIN) go install github.com/gogo/protobuf/$(notdir $@)

PROTOC_GEN_LINT := $(PROTO_GOBIN)/protoc-gen-lint
$(PROTOC_GEN_LINT): $(MODFILE_DIR)/github.com/ckaznocha/protoc-gen-lint/UPDATE_CHECK $(PROTO_GOBIN)
	@echo "+ $@"
	$(SILENT)GOBIN=$(PROTO_GOBIN) go install github.com/ckaznocha/protoc-gen-lint

GOGO_M_STR := Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types

# Hack: there's no straightforward way to escape a comma in a $(subst ...) command, so we have to resort to this little
# trick.
null :=
space := $(null) $(null)
comma := ,

M_ARGS_STR := $(subst $(space),$(comma),$(strip $(ALL_PROTOS_REL)))


$(PROTOC_INCLUDES): $(PROTOC)

GOGO_DIR = $(shell go list -f '{{.Dir}}' -m github.com/gogo/protobuf)

.PHONY: proto-fmt
proto-fmt: $(PROTOC_GEN_LINT)
	@echo "Checking for proto style errors"
	$(SILENT)PATH=$(PROTO_GOBIN) $(PROTOC) \
		-I$(PROTOC_INCLUDES) \
		-I$(GOGO_DIR)/protobuf \
		--lint_out=. \
		--proto_path=$(PROTO_BASE_PATH) \
		$(ALL_PROTOS)

PROTO_DEPS=$(PROTOC) $(PROTOC_INCLUDES)

###############
## Utilities ##
###############

.PHONY: printapisrcs
printapisrcs:
	@echo $(GENERATED_PB_SRCS)

#######################################################################
## Generate gRPC proto messages, services for the API. ##
#######################################################################

# Generate all of the proto messages and gRPC services with one invocation of
# protoc when any of the .pb.go sources don't exist or when any of the .proto
# files change.
$(GENERATED_BASE_PATH)/%.pb.go: $(PROTO_BASE_PATH)/%.proto $(PROTO_DEPS) $(PROTOC_GEN_GO_BIN) $(ALL_PROTOS)
	@echo "+ $@"
	mkdir -p $(dir $@)
	$(SILENT)PATH=$(PROTO_GOBIN) $(PROTOC) \
		-I$(GOGO_DIR) \
		-I$(PROTOC_INCLUDES) \
		--proto_path=$(PROTO_BASE_PATH) \
		--gofast_out=$(GOGO_M_STR:%=%,)$(M_ARGS_STR:%=%,)plugins=grpc:$(GENERATED_BASE_PATH) \
		$(dir $<)/*.proto


# Nukes pretty much everything that goes into building protos.
# You should not have to run this day-to-day, but it occasionally is useful
# to get out of a bad state after a version update.
.PHONY: clean-proto-deps
clean-proto-deps:
	@echo "+ $@"
	rm -f $(PROTOC_FILE)
	rm -rf $(PROTOC_DIR)
	rm -f $(PROTO_GOBIN)
