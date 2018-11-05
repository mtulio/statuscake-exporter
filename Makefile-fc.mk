# ##################
# Makefile functions

define deps_tag
	# @if [[ $(version)x == "x" ]]; then \
	# 	echo -e "\n Error: the version was not specified."; \
	# 	$(call show_usage) \
	# 	exit 1; \
	# fi
	@if [[ "$(message)"x == "x" ]]; then \
		echo -e "\n Error: the commit message was not provided."; \
		$(call show_usage) \
		exit 1; \
	fi
endef