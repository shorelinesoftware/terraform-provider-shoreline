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

install: tf-clean
	rm -f .terraform.lock.hcl 
	make -C .. install 
	$(TERRAFORM) init
