# Style checks
STYLE_CHECK := ./tools/style/entrypoints/check-style.sh
STYLE_VERBOSE_FLAG := $(if $(filter true TRUE yes YES 1,$(VERBOSE)),--verbose,)

.PHONY: style style-verbose style-all style-all-verbose style-all-strict style-all-strict-verbose

##@ Style
style: ## Run required STYLE.md checks (set VERBOSE=true for details)
	$(STYLE_CHECK) --profile required $(STYLE_VERBOSE_FLAG)

style-verbose: ## Run required STYLE.md checks with detailed output
	$(MAKE) style VERBOSE=true

style-all: ## Run required checks and recommendation checks (set VERBOSE=true for details)
	$(STYLE_CHECK) --profile all $(STYLE_VERBOSE_FLAG)

style-all-verbose: ## Run style-all checks with detailed output
	$(MAKE) style-all VERBOSE=true

style-all-strict: ## Run style-all checks and fail on recommendation findings (set VERBOSE=true for details)
	$(STYLE_CHECK) --profile all --strict-recommendations $(STYLE_VERBOSE_FLAG)

style-all-strict-verbose: ## Run style-all-strict checks with detailed output
	$(MAKE) style-all-strict VERBOSE=true
