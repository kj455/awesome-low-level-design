.PHONY: lld

lld:
	@cd scripts/lld && go run . $(filter-out $@,$(MAKECMDGOALS))

%:
	@:
