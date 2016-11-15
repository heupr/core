package conflation

type Scenario2b struct {
  Algorithm Conflation
}

var keywords = []string{"Close #", "Closes #", "Closed #", "Fix #", "Fixes #", "Fixed #", "Resolve #", "Resolves #", "Resolved #"}
// TODO: evaluate optimization
// There could be a better way to handle this logic. Once our unit testing is
// robust @taylor will play around (if needed for performance).
func (c *Scenario2b) extractIssueId(pull *github.PullRequest) int {
	fixIdx := 0
	for i := 0; i < len(keywords); i++ {
		fixIdx = strings.LastIndex(*pull.Body, keywords[i])
		if fixIdx != -1 {
			break
		}
	}
	if fixIdx == -1 {
		return -1
	}
	body := string(*pull.Body)
	body = body[fixIdx:]
	digit := digitRegexp.Find([]byte(body))
	s, _ := strconv.ParseInt(string(digit), 10, 32)
	return int(s)
}


func (s *Scenario2b) IsValid(pull github.PullRequest) bool {
  return true
}

func (s *Scenario2b) Filter(pull github.PullRequest) {
  if IsValid(pull) {
    //append s.Algorithm.Context.Pulls
  }
}
