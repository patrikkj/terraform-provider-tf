### Local development

Rough draft in `./mk download_providers`, needs updates

```hcl
# .terraformrc
provider_installation {

  dev_overrides {
      "local.dev/patrikkj/ssh" = "/Users/pakj/go/bin"
  }

  filesystem_mirror {
    path = "/Users/pakj/.terraform.d/plugins"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

### File structure

```bash
pakj at m4 in ~/.terraform.d/plugins/local.dev/patrikkj
$ tl5
.
└── ssh
    └── 0.0.1
        └── darwin_arm64
            └── terraform-provider-ssh_v0.0.1
```
