module shoreline.io/terraform/terraform-provider-shoreline

go 1.15

require (
	github.com/MichaelMure/go-term-markdown v0.1.4
	github.com/c-bata/go-prompt v0.2.6
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gizak/termui/v3 v3.1.0
	github.com/guptarohit/asciigraph v0.5.2
	github.com/hashicorp/terraform-plugin-docs v0.4.0
	github.com/hashicorp/terraform-plugin-sdk v1.17.2
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.6.1
	github.com/jedib0t/go-pretty/v6 v6.2.2
	github.com/klauspost/compress v1.11.2
	github.com/spf13/viper v1.7.1
	gopkg.in/yaml.v2 v2.3.0
)

replace shoreline.io/terraform/terraform-provider-shoreline/provider => ./provider
