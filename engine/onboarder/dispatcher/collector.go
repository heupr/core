package dispatcher

import "coralreefci/engine/onboarder/retriever"

var Workload = make(chan *retriever.RepoData, 100)

func Collector(repodata map[int]*retriever.RepoData) {
	if len(repodata) != 0 {
		for _, rd := range repodata {
			Workload <- rd
		}
	}
}
