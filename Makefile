BOLD := $(shell tput bold)
RESET := $(shell tput sgr0)

.DEFAULT_GOAL := help

.PHONY: build # Build flux
build:
	@echo "${BOLD}Building flux...${RESET}"
	@cd cli && go build -o ../flux

.PHONY: keep # Copy the flux binary to /usr/local/bin
keep: build 
	@echo "${BOLD}Copying flux to /usr/local/bin...${RESET}"
	@sudo cp flux /usr/local/bin/flux

.PHONY: update # Update software within Dockerfiles to latest stable release
update:
	@echo "${BOLD}Updating Dockerfiles to latest stable releases...${RESET}"
	@python3 update.py

.PHONY: images # Build all the images and remove any dangling ones
images: build
	@echo "${BOLD}Building images...${RESET}"
	@./flux build -t main
	@./flux build -t cuda -f devel/cuda.Dockerfile
	@./flux build -t func -f devel/func.Dockerfile -u root
	@./flux build -t tex -f devel/tex.Dockerfile
	@docker image prune -f

.PHONY: main # Build the main image
main: build
	@echo "${BOLD}Building main image...${RESET}"
	@./flux build -t main
	@docker image prune -f

.PHONY: help # Display the help message
help:
	@echo "${BOLD}Available targets:${RESET}"
	@cat Makefile | grep '.PHONY: [a-z]' | sed 's/.PHONY: / /g' | sed 's/ #* / - /g'
