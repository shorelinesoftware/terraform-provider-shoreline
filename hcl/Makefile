#TERRAFORM=../../terraform
TERRAFORM=terraform

apply:
	$(TERRAFORM) plan
	$(TERRAFORM) apply

tf-clean:
	rm -f terraform.tfstate

install: tf-clean
	rm -f .terraform.lock.hcl 
	make -C .. install 
	$(TERRAFORM) init
