package github

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccGithubGistDataSource(t *testing.T) {
	t.Run("reads gist without error", func(t *testing.T) {
		config := `data "github_gist" "test" {
			gist_id = "2ddc5d178f286766751072ca103e6f3d"
		}`

		check := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet("data.github_gist.test", "id"),
			resource.TestCheckResourceAttrSet("data.github_gist.test", "description"),
			resource.TestCheckResourceAttrSet("data.github_gist.test", "public"),
			resource.TestCheckResourceAttrSet("data.github_gist.test", "owner"),
			resource.TestCheckResourceAttrSet("data.github_gist.test", "file.#"),
			resource.TestCheckResourceAttrSet("data.github_gist.test", "html_url"),
			resource.TestCheckResourceAttrSet("data.github_gist.test", "created_at"),
			resource.TestCheckResourceAttrSet("data.github_gist.test", "updated_at"),
			resource.TestCheckResourceAttrSet("data.github_gist.test", "node_id"),
		)

		testCase := func(t *testing.T, mode string) {
			resource.Test(t, resource.TestCase{
				PreCheck:  func() { skipUnlessMode(t, mode) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: config,
						Check:  check,
					},
				},
			})
		}

		t.Run("with an anonymous account", func(t *testing.T) {
			testCase(t, anonymous)
		})

		t.Run("with an individual account", func(t *testing.T) {
			testCase(t, individual)
		})

		t.Run("with an organization account", func(t *testing.T) {
			t.Skip("organization account not supported for this operation")
		})
	})
}
