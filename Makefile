BOLD := \033[1m
RESET := \033[0m

.DEFAULT_GOAL := help

.PHONY: flux # Build flux
flux:
	@echo "${BOLD}Building flux...${RESET}"
	@cp README.md cli/info/README.md
	@cd cli && go build -o ../flux

.PHONY: keep # Copy the flux binary to /usr/local/bin
keep: flux
	@echo "${BOLD}Copying flux to /usr/local/bin...${RESET}"
	@sudo cp flux /usr/local/bin/flux

.PHONY: update # Update software within Dockerfiles to latest stable release
update:
	@echo "${BOLD}Updating Dockerfiles to latest stable releases...${RESET}"
	@python3 update.py

.PHONY: images # Build all the images and remove any dangling ones
images: flux
	@echo "${BOLD}Building images...${RESET}"
	@./flux build -t main
	@./flux build -t cuda -f devel/cuda.Dockerfile
	@./flux build -t func -f devel/func.Dockerfile -u root
	@./flux build -t tex -f devel/tex.Dockerfile
	@docker image prune -f

.PHONY: main # Build the main image
main: flux
	@echo "${BOLD}Building main image...${RESET}"
	@./flux build -t main
	@docker image prune -f

.PHONY: help # Display the help message
help:
	@echo "${BOLD}Available targets:${RESET}"
	@cat Makefile | grep '.PHONY: [a-z]' | sed 's/.PHONY: / /g' | sed 's/ #* / - /g'
