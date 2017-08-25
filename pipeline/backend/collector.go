package backend

var Workload = make(chan *RepoData, 100)

func Collector(repodata map[int]*RepoData) {
	if len(repodata) != 0 {
		for _, rd := range repodata {
			Workload <- rd
		}
	}
}
