---
layout: "github"
page_title: "GitHub: github_gist"
description: |-
  Get information on an existing Github gist
---

# github_gist

Use this data source to retrieve information about a GitHub gist.

## Example Usage

```hcl
data "github_gist" "example" {
  gist_id = "3a641566d0a49014cea6f395e8be68c42b9e46d9"
}
```

## Argument Reference

- `gist_id` - (Required) Id of the gist.

## Attributes Reference

- `description` - 
- `public` - 
- `owner` - 
- `files` - 
- `filename` - 
- `language` - 
- `raw_url` - 
- `size` - 
- `type` - 
- `html_url` - 
- `created_at` - 
- `updated_at` - 
- `node_id` - 