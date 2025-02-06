# terraform-provider-{{ .ProviderShortName }}

This repository contains the terraform provider implementation and docs for the {{ .RenderedProviderName }} APIs.

## Documentation

Documentation located in the `/docs` directory is compiled and generated using the Terraform docs plugin based on the local `templates/**/*.md.tmpl` files.


### Running acceptance tests

Acceptance tests may be run against a local deployment of {{ .RenderedProviderName }}. In order for these tests to work, the provider devcontainer needs to run in the same `shoreline-net` podman network as the other podman containers related to {{ .RenderedProviderName }}. This will allow the provider to have access to the ceph gateway in order to upload `{{ .ProviderShortName }}_file`s resources to the local S3 deployment.

To do that, simply uncomment this line in `devcontainer.json`:
```
    // "runArgs": ["--network=shoreline-net"],
```
and then rebuild the container, making sure the podman network was created by the shoreline-in-a-box run script.


### Troubleshooting / FAQ

###### Creating `{{ .ProviderShortName }}_file`s via terraform does not work
Check backend container for error logs.

Ensure you've set the `ENABLE_LOCAL_OP_COPY` flag (see above for more information).

Also, ensure that the AWS credentials are set on the ceph containers and on the backend (`aws sso login` then `podman login` with the default AWS credentials). Credentials expire after some time (~12 hours) which may cause op copy to stop working all of a sudden. This is because the backend will fail to presign the PUT url if the AWS credentials are invalid / expired / not set.