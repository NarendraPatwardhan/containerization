BOLD := \033[1m
RESET := \033[0m

.DEFAULT_GOAL := help

.PHONY: build # Build flux
build:
	@echo "${BOLD}Building flux...${RESET}"
	@cd cli && go build -o ../flux

.PHONY: help # Display the help message
help:
	@echo "${BOLD}Available targets:${RESET}"
	@cat Makefile | grep '.PHONY: [a-z]' | sed 's/.PHONY: / /g' | sed 's/ #* / - /g'