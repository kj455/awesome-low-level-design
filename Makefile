.PHONY: lld

lld:
	@bash scripts/lld.sh $(filter-out $@,$(MAKECMDGOALS))

%:
	@:
