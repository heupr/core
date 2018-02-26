package backend

import (
	"testing"

	"core/pipeline/gateway/conflation"
)

func TestDispatcher(t *testing.T) {
	repoID := int64(5555)
	s := &Server{
		Repos: &ActiveRepos{
			Actives: map[int64]*ArchRepo{
				repoID: &ArchRepo{
					Hive: &ArchHive{
						Blender: &Blender{
							Conflator: &conflation.Conflator{
								Context: &conflation.Context{
									Issues: []conflation.ExpandedIssue{
										conflation.ExpandedIssue{},
									},
								},
								Normalizer: conflation.Normalizer{
									Context: &conflation.Context{
										Issues: []conflation.ExpandedIssue{
											conflation.ExpandedIssue{},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	workload <- &RepoData{
		RepoID: repoID,
	}
	s.Dispatcher(1)
}
