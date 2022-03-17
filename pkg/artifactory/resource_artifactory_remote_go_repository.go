package artifactory

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type GoRemoteRepo struct {
	RemoteRepositoryBaseParams
	VcsType        string `json:"vcsType"`
	VcsGitProvider string `json:"vcsGitProvider"`
}

func resourceArtifactoryRemoteGoRepository() *schema.Resource {
	const packageType = "go"

	var goRemoteSchema = mergeSchema(baseRemoteRepoSchema, map[string]*schema.Schema{
		"vcs_type": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "GIT",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"GIT"}, false)),
			Description:      `(Optional) Artifactory supports proxying the Git providers. Default value is "GIT".`,
		},
		"vcs_git_provider": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "ARTIFACTORY",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"GITHUB", "ARTIFACTORY"}, false)),
			Description:      `(Optional) Artifactory supports proxying the following Git providers out-of-the-box: GitHub or a remote Artifactory instance. Default value is "ARTIFACTORY".`,
		},
	}, repoLayoutRefSchema("remote", packageType))

	var unpackGoRemoteRepo = func(s *schema.ResourceData) (interface{}, string, error) {
		d := &ResourceData{s}
		repo := GoRemoteRepo{
			RemoteRepositoryBaseParams: unpackBaseRemoteRepo(s, packageType),
			VcsType:                    d.getString("vcs_type", false),
			VcsGitProvider:             d.getString("vcs_git_provider", false),
		}
		return repo, repo.Id(), nil
	}

	return mkResourceSchema(goRemoteSchema, defaultPacker, unpackGoRemoteRepo, func() interface{} {
		repoLayout, _ := getDefaultRepoLayoutRef("remote", packageType)()
		return &GoRemoteRepo{
			RemoteRepositoryBaseParams: RemoteRepositoryBaseParams{
				Rclass:              "remote",
				PackageType:         packageType,
				RemoteRepoLayoutRef: repoLayout.(string),
			},
		}
	})
}
