
# Pray that nothing breaks...
update_deps:
	@rm -rf vendor/
	@echo "--> Running dep"
	@dep ensure
