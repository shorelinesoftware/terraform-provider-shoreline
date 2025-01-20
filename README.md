# terraform-provider-shoreline

This repository contains the terraform provider implementation and docs for the Shoreline Software APIs.

## Documentation

Documentation located in the `/docs` directory is compiled and generated from a handful of disparate sources:

- The [Terraform docs plugin](https://github.com/hashicorp/terraform-plugin-docs) uses local `templates/**/*.md.tmpl` files to compile Terraform resource files.
- `content/**/*.md` files are the base templates used by the `tfdocsplugin`, but they optionally use terminology and/or relative path URLs based on the docs.shoreline.io content system.

  The `content/terms.json` file defines a series of terminology patterns that map to Docs URLs. These same terminology paths are usable within `content/**/*.md` files using the same `/t/<term>` URL syntax. See [Docs: Terminology Links](https://docs.shoreline.io/internal/writing#terminology-links) for more details.

  For example, consider the following template Markdown:

  ```md
  Actions execute shell commands on associated [Resources](/t/resource). Whenever an [Alarm](/t/alarm) fires the associated [Bot](/t/bot) triggers the corresponding [Action](/t/action), closing the basic auto-remediation loop of Shoreline.

  - [name](/actions/properties#name) - The name of the Action.
  ```

  The final Markdown is converted to the following, automatically mapping terminology links to the appropriate external docs.shoreline.io Article link.

  ```md
  Actions execute shell commands on associated [Resources](https://docs.shoreline.io/platform/resources). Whenever an [Alarm](https://docs.shoreline.io/alarms) fires the associated [Bot](https://docs.shoreline.io/bots) triggers the corresponding [Action](https://docs.shoreline.io/actions), closing the basic auto-remediation loop of Shoreline.

  - [name](https://docs.shoreline.io/actions/properties#name) - The name of the Action.
  ```

### Build the Documentation

1. _(Optional)_ Install all required Node modules with `yarn install`.
2. Build the templates and generate the final docs with `gulp build` (or just `gulp`):

   ```
   [15:26:28] Using gulpfile F:\projects\shoreline\repos\terraform\terraform-provider-shoreline\gulpfile.js
   [15:26:28] Starting 'build'...
   [15:26:28] Starting 'buildTemplates'...
   [15:26:28] Finished 'buildTemplates' after 46 ms
   [15:26:28] Starting 'generateDocs'...
   rendering website for provider "terraform-provider-shoreline"
   copying any existing content to tmp dir
   exporting schema from Terraform
   compiling provider "shoreline"
   generating missing resource content
   resource "shoreline_action" template exists, skipping
   generating template for "shoreline_alarm"
   generating template for "shoreline_bot"
   generating template for "shoreline_file"
   generating template for "shoreline_metric"
   generating template for "shoreline_resource"
   generating missing data source content
   generating missing provider content
   provider "terraform-provider-shoreline" template exists, skipping
   cleaning rendered website dir
   rendering templated website to static markdown
   rendering "index.md.tmpl"
   rendering "resources\\action.md.tmpl"
   rendering "resources\\alarm.md.tmpl"
   rendering "resources\\bot.md.tmpl"
   rendering "resources\\file.md.tmpl"
   rendering "resources\\metric.md.tmpl"
   rendering "resources\\resource.md.tmpl"
   [15:26:41] Finished 'generateDocs' after 13 s
   [15:26:41] Finished 'build' after 13 s
   ```

   This process builds `templates/**/*.md.tmpl` files from `content/**/*.md` files, replacing any terminology/relative path URLs, then generates the `docs/` files via `tfdocsplugin`.


### Running acceptance tests

Acceptance tests may be run against a local deployment of shoreline. In order for these tests to work, the provider devcontainer needs to run in the same `shoreline-net` podman network as the other podman containers related to Shoreline. This will allow the provider to have access to the ceph gateway in order to upload `shoreline_file`s resources to the local S3 deployment.

To do that, simply uncomment this line in `devcontainer.json`:
```
    // "runArgs": ["--network=shoreline-net"],
```
and then rebuild the container, making sure the podman network was created by the shoreline-in-a-box run script.


### Troubleshooting / FAQ

###### Creating `shoreline_file`s via terraform does not work
Check backend container for error logs.

Ensure you've set the `ENABLE_LOCAL_OP_COPY` flag (see above for more information).

Also, ensure that the AWS credentials are set on the ceph containers and on the backend (`aws sso login` then `podman login` with the default AWS credentials). Credentials expire after some time (~12 hours) which may cause op copy to stop working all of a sudden. This is because the backend will fail to presign the PUT url if the AWS credentials are invalid / expired / not set.