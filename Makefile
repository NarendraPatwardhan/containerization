BOLD := \033[1m
RESET := \033[0m

.DEFAULT_GOAL := help

.PHONY: flux # Build flux
flux:
	@echo "${BOLD}Building flux...${RESET}"
	@cp README.md cli/info/README.md
	@cd cli && go build -o ../flux

.PHONY: images # Build all the images
images: flux
	@echo "${BOLD}Building images...${RESET}"
	@./flux build -t main
	@./flux build -t cuda -f devel/cuda.Dockerfile
	@./flux build -t func -f devel/func.Dockerfile -u root
	@./flux build -t tex -f devel/tex.Dockerfile

.PHONY: help # Display the help message
help:
	@echo "${BOLD}Available targets:${RESET}"
	@cat Makefile | grep '.PHONY: [a-z]' | sed 's/.PHONY: / /g' | sed 's/ #* / - /g'