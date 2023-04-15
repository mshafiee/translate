# Define variables for the output directories and executable name
OUTPUT_DIR := ./bin
EXEC_NAME := translate

# Define the operating systems and architectures to build for
OS_LIST := darwin linux windows
ARCH_LIST := amd64 arm64

# Generate a list of build targets using the cross-platform build syntax
	TARGETS := $(foreach os,$(OS_LIST),$(foreach arch,$(ARCH_LIST),$(OUTPUT_DIR)/$(os)-$(arch)/$(EXEC_NAME)))

# Set the default build target to build for all platforms
all: $(TARGETS)

# Define the build rules for each platform
$(OUTPUT_DIR)/%/$(EXEC_NAME): cmd/main.go
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) go build -o $@ $<

# Define a clean rule to remove all output directories
clean:
	rm -rf $(OUTPUT_DIR)

# Ensure the output directory exists before building
$(TARGETS): | $(OUTPUT_DIR)

# Define a phony target for the output directory to avoid conflicts with files
.PHONY: $(OUTPUT_DIR)

# Create the output directory if it does not exist
$(OUTPUT_DIR):
	mkdir -p $(foreach os,$(OS_LIST),$(foreach arch,$(ARCH_LIST),$(OUTPUT_DIR)/$(os)-$(arch)))
