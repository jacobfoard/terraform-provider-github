package github

import (
	"fmt"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceGithubGists() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGithubGistsRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"since": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRFC3339TimeString,
			},
			"gists": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gistDataSchemaResource(false),
			},
		},
	}
}

func dataSourceGithubGistsRead(d *schema.ResourceData, meta interface{}) error {
	owner := meta.(*Owner)

	username := owner.name

	if providedUsername, ok := d.GetOk("username"); ok {
		username = providedUsername.(string)
	}

	opt := &github.GistListOptions{
		ListOptions: github.ListOptions{
			PerPage: maxPerPage,
		},
	}

	sinceRaw, ok := d.GetOk("since")
	if ok {
		since, err := time.Parse(time.RFC3339, sinceRaw.(string))
		if err != nil {
			return err
		}

		opt.Since = since
	}

	gists := make([]interface{}, 0)

	for {
		gistList, resp, err := owner.v3client.Gists.List(owner.StopContext, username, opt)
		if err != nil {
			return err
		}

		result, err := flattentGists(gistList)
		if err != nil {
			return err
		}

		gists = append(gists, result...)

		if resp.NextPage == 0 {
			break
		}

		opt.ListOptions.Page = resp.NextPage
	}

	d.SetId(fmt.Sprintf("%s-gists", username))
	if err := d.Set("gists", gists); err != nil {
		return fmt.Errorf("error setting gists: %s", err)
	}

	return nil
}

func flattentGists(gists []*github.Gist) ([]interface{}, error) {
	if len(gists) == 0 {
		return nil, nil
	}

	flatGists := make([]interface{}, len(gists))

	for i, gist := range gists {
		flatGist, err := flattenGist(gist)
		if err != nil {
			return nil, err
		}

		flatGists[i] = flatGist
	}

	return flatGists, nil
}
