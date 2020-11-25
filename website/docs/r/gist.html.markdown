---
layout: "github"
page_title: "GitHub: github_gist"
description: |-
  Creates and manages GitHub gists.
---

# github_gist

This resource allows you to create and manage Github gists.

Note that the provider `token` must have the gist scope attached to it.

## Example Usage

```hcl
resource "github_gist" "readme" {
  description = "example"
  file {
      name = "README.md"
      content = <<-EOF
      # terraform-github-provider
      EOF
  }
}
```

## Argument Reference

The following arguments are supported:

- `file` - (Required) See `file` block below

- `description` - (Optional) Addition info about the gist.

- `public` - (Optional) When set to true the gist will be public. When set to false it is still publicly accessible, but not shown on the user's profile. Defaults to `false`

---

The `file` block consists of:

- `name` - (Required) Full name of the file.
- `content` - (Required) String representing the content of the file.

## Attribute Reference

The following additional attributes are exported:

- `gist_id` - Id for referencing the gist.

- `html_url` - URL for the gist in the standard acess ui.

- `git_pull_url` - URL to use for pulling the underlying repo of the gist.

- `git_push_url` - URL to use for pushing to the underlying repo of the gist.

- `created_at` - Date of gist creation.

- `updated_at` - Date of most recent gist update.

- `node_id` - The Node Id for use with the GraphQL API.

- `file[0].raw_url` - URL for the "raw" view of the file

## Import

GitHub Branch can be imported using an ID made up of `repository:branch`, e.g.

```
$ terraform import github_branch.terraform terraform:master
```

Optionally, a source branch may be specified using an ID of `repository:branch:source_branch`.
This is useful for importing branches that do not branch directly off master.

```
$ terraform import github_branch.terraform terraform:feature-branch:dev
```
