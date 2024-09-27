# Get the directory where the Makefile is located
MAKEFILE_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# Load the .env file and export variables
env_file=$(MAKEFILE_DIR)/.env
include $(env_file)
export $(env_file)

# Variables
date_stamp=$(shell date +%Y%m%d%H%M%S)
img_name=$(FQDN_REGISTRY)/$(STACK_NAME)

# Printing vars for debugging
# $(info MAKEFILE_DIR = $(MAKEFILE_DIR))
# $(info ROOT_DIR = $(ROOT_DIR))
# $(info FQDN_REGISTRY = $(FQDN_REGISTRY))
# $(info img_name = $(img_name))

# Build production image & push to repo
build:
	docker build -f $(ROOT_DIR)/volumes/app/Dockerfile -t $(img_name):$(date_stamp) $(ROOT_DIR)/volumes/app
	docker push $(img_name):$(date_stamp)

# Build development image & push to repo
build_dev:
	docker build -f $(ROOT_DIR)/volumes/app/Dockerfile.dev -t $(img_name):dev $(ROOT_DIR)/volumes/app
	docker push $(img_name):dev
