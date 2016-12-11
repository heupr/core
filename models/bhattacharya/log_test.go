package bhattacharya

import "testing"

// DOC: This is the unit test for the logging output structure; the current
//      setting is for the logs/ directory.
func TestLog(t *testing.T) {
	testlog := CreateLog("log-test", false)
	testlog.Log("Testing input values")
	testlog.Flush()
}
