#TERRAFORM=../../terraform
TERRAFORM=terraform

plan:
	$(TERRAFORM) plan

apply:
	$(TERRAFORM) apply

destroy:
	$(TERRAFORM) destroy

tf-clean:
	rm -f terraform.tfstate
	rm -f .terraform.lock.hcl
	rm -rf .terraform

install:
	rm -f .terraform.lock.hcl
	make -C .. install 
	$(TERRAFORM) init

schema:
	@terraform providers schema -json
