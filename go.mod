module shoreline.io/terraform/terraform-provider-shoreline

go 1.16

require (
	github.com/aws/aws-sdk-go v1.37.0 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/hashicorp/go-getter v1.5.11 // indirect
	github.com/hashicorp/hcl/v2 v2.8.2 // indirect
	github.com/hashicorp/terraform-json v0.13.0 // indirect
	github.com/hashicorp/terraform-plugin-docs v0.8.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.7.0
	github.com/klauspost/compress v1.11.2
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/viper v1.7.1
	golang.org/x/tools v0.0.0-20201028111035-eafbe7b904eb // indirect
	google.golang.org/api v0.34.0 // indirect
)

replace shoreline.io/terraform/terraform-provider-shoreline/provider => ./provider
