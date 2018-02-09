package backend

var workload = make(chan *RepoData, 100)

func collector(repodata map[int64]*RepoData) {
	if len(repodata) != 0 {
		for _, rd := range repodata {
			workload <- rd
		}
	}
}
