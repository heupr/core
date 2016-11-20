package conflation

import (
	"github.com/google/go-github/github"
	"testing"
)

func TestConflater(t *testing.T) {
	context := &Context{}

	comboScenarios := []Scenario{&Scenario2b{}}
	conflationAlgorithms := []ConflationAlgorithm{&ComboAlgorithm{Context: context}}
	conflator := Conflator{Scenarios: comboScenarios, ConflationAlgorithms: conflationAlgorithms, Context: context}

	number := 12886
	issue := "The following are not implementing ICloneable but according to Net Standard 2.0 they should:\r\n\r\n- [x] System.Array.ArrayEnumerator\r\n- [x] System.Array.SZArrayEnumerator\r\n- [x] System.Collections.ArrayList\r\n- [x] System.Collections.ArrayList.ArrayListEnumerator\r\n- [x] System.Collections.ArrayList.ArrayListEnumeratorSimple\r\n- [x] System.Collections.ArrayList.IListWrapper.IListWrapperEnumWrapper\r\n- [x] System.Collections.BitArray\r\n- [x] System.Collections.BitArray.BitArrayEnumeratorSimple\r\n- [x] System.Collections.Hashtable\r\n- [x] System.Collections.Hashtable.HashtableEnumerator\r\n- [x] System.Collections.Queue\r\n- [x] System.Collections.Queue.QueueEnumerator\r\n- [x] System.Collections.SortedList\r\n- [x] System.Collections.SortedList.SortedListEnumerator\r\n- [x] System.Collections.Stack\r\n- [x] System.Collections.Stack.StackEnumerator\r\n- [x] System.ComponentModel.MaskedTextProvider\r\n- [ ] System.Configuration.Assemblies.AssemblyHash\r\n- [x] System.Delegate\r\n- [x] System.Net.Http.Headers.AuthenticationHeaderValue\r\n- [x] System.Net.Http.Headers.CacheControlHeaderValue\r\n- [x] System.Net.Http.Headers.ContentDispositionHeaderValue\r\n- [x] System.Net.Http.Headers.ContentRangeHeaderValue\r\n- [x] System.Net.Http.Headers.EntityTagHeaderValue\r\n- [x] System.Net.Http.Headers.MediaTypeHeaderValue\r\n- [x] System.Net.Http.Headers.MediaTypeWithQualityHeaderValue\r\n- [x] System.Net.Http.Headers.NameValueHeaderValue\r\n- [x] System.Net.Http.Headers.NameValueWithParametersHeaderValue\r\n- [x] System.Net.Http.Headers.ProductHeaderValue\r\n- [x] System.Net.Http.Headers.ProductInfoHeaderValue\r\n- [x] System.Net.Http.Headers.RangeConditionHeaderValue\r\n- [x] System.Net.Http.Headers.RangeHeaderValue\r\n- [x] System.Net.Http.Headers.RangeItemHeaderValue\r\n- [x] System.Net.Http.Headers.RetryConditionHeaderValue\r\n- [x] System.Net.Http.Headers.StringWithQualityHeaderValue\r\n- [x] System.Net.Http.Headers.TransferCodingHeaderValue\r\n- [x] System.Net.Http.Headers.TransferCodingWithQualityHeaderValue\r\n- [x] System.Net.Http.Headers.ViaHeaderValue\r\n- [x] System.Net.Http.Headers.WarningHeaderValue\r\n- [x] System.OperatingSystem\r\n- [ ] System.Runtime.Remoting.Messaging.CallContextRemotingData\r\n- [ ] System.Runtime.Remoting.Messaging.CallContextSecurityData\r\n- [ ] System.Runtime.Remoting.Messaging.LogicalCallContext\r\n- [x] System.Runtime.Serialization.Formatters.Binary.IntSizedArray\r\n- [x] System.Runtime.Serialization.Formatters.Binary.SizedArray\r\n- [x] System.RuntimeType\r\n- [x] System.Version\r\n- [x] System.Xml.Schema.XmlAtomicValue\r\n- [x] System.Xml.XPath.XPathNavigator\r\n- [x] System.Xml.XPath.XPathNodeIterator\r\n- [x] System.Xml.XmlNode\r\n- [x] string"
	title := "Issue Title"
	iAssignee := "Mike"
	issueAssignee := github.User{Name: &iAssignee}
	gitIssue := github.Issue{Number: &number, Title: &title, Body: &issue, Assignee: &issueAssignee}
	issues := []github.Issue{gitIssue}

	pullNumber := 2
	pull := "If NetworkAddressChanged has a subscriber and NetworkAvailabilityChanged doesn't, attempting to remove a delegate from NetworkAvailabilityChanged results in a NullReferenceException as a null timer static field is dereferenced.\r\n\r\nThis was failing intermittently in tests because the NetworkAddressChanged and NetworkAvailabilityChanged tests were in different classes and thus could potentially run in parallel.  If the NetworkAvailabilityChanged_JustRemove_Success test ran while one of the NetworkAddressChanged tests was in flight, it would fail with a null ref.\r\n\r\nI've fixed the product bug by adding a null check.  I've also updated the tests so that all of these event handler tests are in the same class, and added explicit tests for this case that will deterministically fail if the product regresses.\r\n\r\nFixes https://github.com/dotnet/corefx/issues/12886\r\ncc: @mellinoe, @davidsh, @cipop"
	pullTitle := "Pull Title"
	pAssignee := "John"
	pullAssignee := github.User{Name: &pAssignee}
	issueURL := "https://github.com/dotnet/corefx/issues/12886"
	gitPull := github.PullRequest{Number: &pullNumber, Title: &pullTitle, Body: &pull, IssueURL: &issueURL, Assignee: &pullAssignee}
	pulls := []github.PullRequest{gitPull}

	conflator.Context.Issues = make([]ExpandedIssue, len(issues)+len(pulls))
	conflator.SetIssueRequests(issues)
	conflator.SetPullRequests(pulls)

	conflator.Conflate()

	//TODO: fix assert
	/*
	  for i := 0; i < len(conflator.Context.Issues); i++ {
			if conflator.Context.Issues[i].CrIssueType.GitIssueType.Number == 12886 {
					Assert(t, *pullAssignee.Name, *conflator.Context.Issues[i].CrIssueType.GitIssueType.Assignee.Name, "Mike/John Issue/Pull")
			}
		}*/
}

func Assert(t *testing.T, expected string, actual string, input string) {
	if actual != expected {
		t.Error(
			"\nFOR:       ", input,
			"\nEXPECTED:  ", expected,
			"\nACTUAL:    ", actual,
		)
	}
}
