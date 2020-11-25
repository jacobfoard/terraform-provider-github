package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGithubGist() *schema.Resource {
	return &schema.Resource{
		Create: resourceGithubGistCreate,
		Update: resourceGithubGistUpdate,
		Read:   resourceGithubGistRead,
		Delete: resourceGithubGistDelete,

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"file": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"raw_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"gist_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"html_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_pull_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_push_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGithubGistCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Owner).v3client
	ctx := context.Background()

	description := d.Get("description").(string)
	public := d.Get("public").(bool)

	files, err := expandGistFiles(d.Get("file").([]interface{}))
	if err != nil {
		return nil
	}

	gist := &github.Gist{
		Description: &description,
		Public:      &public,
		Files:       files,
	}

	resp, _, err := client.Gists.Create(ctx, gist)
	if err != nil {
		return err
	}

	id := *resp.ID

	d.SetId(id)
	d.Set("gist_id", id)

	return resourceGithubGistRead(d, meta)
}

func resourceGithubGistUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Owner).v3client
	ctx := context.Background()

	description := d.Get("description").(string)
	public := d.Get("public").(bool)

	gist, _, err := client.Gists.Get(ctx, d.Id())
	if err != nil {
		return err
	}

	updatedFileMap, err := expandGistFiles(d.Get("file").([]interface{}))
	if err != nil {
		return err
	}

	for key := range gist.Files {
		if _, ok := updatedFileMap[key]; !ok {
			updatedFileMap[key] = github.GistFile{Filename: nil}
			continue
		}
	}

	updatedGist := &github.Gist{
		Description: &description,
		Public:      &public,
		Files:       updatedFileMap,
	}

	_, _, err = client.Gists.Edit(ctx, d.Id(), updatedGist)
	if err != nil {
		return err
	}

	return resourceGithubGistRead(d, meta)
}

func resourceGithubGistRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Owner).v3client
	ctx := context.Background()

	gist, _, err := client.Gists.Get(ctx, d.Id())
	if err != nil {
		return err
	}

	files := []interface{}{}

	for _, val := range gist.Files {
		file := make(map[string]interface{})

		file["name"] = val.GetFilename()
		file["content"] = val.GetContent()
		file["raw_url"] = val.GetRawURL()

		files = append(files, file)
	}

	d.Set("file", files)
	d.Set("html_url", gist.GetHTMLURL())
	d.Set("git_pull_url", gist.GetGitPullURL())
	d.Set("git_push_url", gist.GetGitPushURL())
	d.Set("created_at", gist.GetCreatedAt().String())
	d.Set("updated_at", gist.GetCreatedAt().String())
	d.Set("node_id", gist.GetNodeID())

	return nil
}

func resourceGithubGistDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Owner).v3client
	ctx := context.Background()

	_, err := client.Gists.Delete(ctx, d.Id())

	return err
}

func expandGistFiles(fileSchemas []interface{}) (map[github.GistFilename]github.GistFile, error) {
	files := make(map[github.GistFilename]github.GistFile)

	for i, val := range fileSchemas {
		fileMap := val.(map[string]interface{})

		var (
			filename, content string
			ok                bool
		)

		filename, ok = fileMap["name"].(string)
		if !ok {
			return nil, fmt.Errorf("unable to find \"name\" in file [%d]", i)
		}

		content, ok = fileMap["content"].(string)
		if !ok {
			return nil, fmt.Errorf("unable to find \"content\" in file [%d]", i)
		}

		files[github.GistFilename(filename)] = github.GistFile{
			Filename: &filename,
			Content:  &content,
		}
	}

	return files, nil
}
