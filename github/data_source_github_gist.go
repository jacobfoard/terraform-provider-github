package github

import (
	"context"

	"github.com/google/go-github/v32/github"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGithubGist() *schema.Resource {
	resource := &schema.Resource{
		Read: dataSourceGithubGistRead,
		Schema: map[string]*schema.Schema{
			"gist_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}

	for key, val := range gistDataSchemaResource(true).Schema {
		resource.Schema[key] = val
	}
	return resource
}

func dataSourceGithubGistRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Owner).v3client

	id := d.Get("gist_id").(string)

	gist, _, err := client.Gists.Get(context.Background(), id)
	if err != nil {
		return err
	}

	result, err := flattenGist(gist)
	if err != nil {
		return err
	}

	d.SetId(id)

	for key, val := range result.(map[string]interface{}) {
		d.Set(key, val)
	}

	return nil
}

func gistDataSchemaResource(includeContent bool) *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"files": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filename": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"language": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"raw_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"html_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	if includeContent {
		resource.Schema["files"].Elem.(*schema.Resource).Schema["content"] = &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		}
	}
	return resource
}

func flattenGist(gist *github.Gist) (interface{}, error) {
	if gist == nil {
		return nil, nil
	}

	result := make(map[string]interface{})

	result["id"] = gist.GetID()
	result["description"] = gist.GetDescription()
	result["public"] = gist.GetPublic()
	result["owner"] = gist.GetOwner().GetLogin()
	result["html_url"] = gist.GetHTMLURL()
	result["created_at"] = gist.GetCreatedAt().String()
	result["updated_at"] = gist.GetUpdatedAt().String()
	result["node_id"] = gist.GetNodeID()

	files := make([]interface{}, len(gist.Files))
	i := 0
	for _, file := range gist.Files {
		fileDetails := make(map[string]interface{})
		fileDetails["filename"] = file.GetFilename()
		fileDetails["language"] = file.GetLanguage()
		fileDetails["raw_url"] = file.GetRawURL()
		fileDetails["size"] = file.GetSize()
		fileDetails["type"] = file.GetType()

		if file.GetContent() != "" {
			fileDetails["content"] = file.GetContent()
		}

		files[i] = fileDetails
		i++
	}

	result["files"] = files

	return result, nil
}
