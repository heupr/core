package conflation

// import "fmt"    // TEMPORARY

type ScenarioAND struct {
	Scenarios []Scenario
}

// DOC: ScenarioAND allows for "bucketing" multiple scenario filters together
//      and applying "AND" logical operations between them (as opposed to the
//      "OR" logic built in to the Conflator).
//      This is a "meta" scenario.
func (s *ScenarioAND) Filter(expandedIssue *ExpandedIssue) bool {
	// for i := 0; i < len(s.Scenarios); i ++ {
	//     if !s.Scenarios[i].Filter(expandedIssue) {
	//         return false
	//     }
	// }
	for _, scenario := range s.Scenarios {
		if !scenario.Filter(expandedIssue) {
			return false
			// fmt.Println("FALSE")    // TEMPORARY
		}
	}
	// fmt.Println("TRUE") // TEMPORARY
	return true
}
