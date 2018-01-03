package labelmaker

import (
	"log"
	"os"
	"testing"

	conf "core/pipeline/gateway/conflation"

	language "cloud.google.com/go/language/apiv1"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
)

var client *language.Client
var ctx context.Context

func setup() {
	ctx = context.Background()

	// Creates a client.
	var err error
	client, err = language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
}

func TestMain(m *testing.M) {
	// setup
	setup()

	retCode := m.Run()

	//teardown()

	// call with result of m.Run()
	os.Exit(retCode)
}

func TestBasicLabelAllocator(t *testing.T) {
	lbModel := LBModel{Classifier: &LBClassifier{Ctx: ctx, Client: client}}
	lbModel.Learn([]string{"A-allocators", "Os-Linux"})
	text1 := `Jemalloc 4.5.0 (first included in 1.21.0) immediately aborts when run on ARM iOS devices. (Upstream issue). Fortunately, the underlying issue seems to have been fixed in jemalloc 5.0.0 and onwards. (Unfortunately, a Rust PR to upgrade to 5.0.1 just got closed: #45163).
Using the system allocator is an okay workaround for now (though, it requires using nightly).`
	Case(t, "A-allocators", text1, lbModel)
	text2 := "System.alloc returns unaligned pointer if align > size"
	Case(t, "A-allocators", text2, lbModel)
}

func Case(t *testing.T, expected string, input string, lbModel LBModel) {
	labels, _ := lbModel.Predict(conf.ExpandedIssue{Issue: conf.CRIssue{github.Issue{Title:github.String(input)}, []int{}, []conf.CRPullRequest{}, github.Bool(false)}})
	if len(labels) == 0 {
		t.Error("INCORRECT LABEL. 0 LABELS RETURNED", "EXPECTING:", expected, "INPUT:", input)
	}
	for i := 0; i < len(labels); i++ {
		if labels[i] != expected {
			t.Error("INCORRECT LABEL", "EXPECTED:", expected, "GOT:", labels[i])
			break
		}
	}
}

func TestBasicYarnIssues(t *testing.T) {
	lbModel := LBModel{Classifier: &LBClassifier{Ctx: ctx, Client: client}}
	lbModel.Learn([]string{"cat-feature", "cat-compatibility", "cat-documentation", "help wanted", "high-priority", "needs-repro-script", "triaged"})
	text := "Request feature for yarn t to act like npm t --shortcut for yarn test"
	Case(t, "cat-feature", text, lbModel)
}


func TestBasicLabelWindows(t *testing.T) {
	lbModel := LBModel{Classifier: &LBClassifier{Ctx: ctx, Client: client}}
	lbModel.Learn([]string{"Os-Windows", "Os-Linux"})

	text := "I tried installing yarn using the installation script from https://yarnpkg.com/en/docs/install#alternatives-tab (I'm on Windows, but I don't have admin rights, so I can't use the Windows installer)."
	labels, _ := lbModel.Predict(conf.ExpandedIssue{Issue: conf.CRIssue{github.Issue{Title: github.String(text)}, []int{}, []conf.CRPullRequest{}, github.Bool(false)}})
	if len(labels) == 0 {
		t.Error("INCORRECT LABEL. 0 LABELS RETURNED", "EXPECTING Os-Windows")
	}
	for i := 0; i < len(labels); i++ {
		if labels[i] != "Os-Windows" {
			t.Error("INCORRECT LABEL", labels[i])
			break
		}
	}
}

func TestYarnWindows(t *testing.T) {
	lbModel := LBModel{Classifier: &LBClassifier{Ctx: ctx, Client: client}}
	lbModel.Learn([]string{"os-windows", "os-linux"})

	text := "Document that Yarn currently doesn't work on Bash on Windows"
	labels, _ := lbModel.Predict(conf.ExpandedIssue{Issue: conf.CRIssue{github.Issue{Title: github.String(text)}, []int{}, []conf.CRPullRequest{}, github.Bool(false)}})
	if len(labels) == 0 {
		t.Error("INCORRECT LABEL. 0 LABELS RETURNED", "EXPECTING Os-Windows")
	}
	for i := 0; i < len(labels); i++ {
		if labels[i] != "Os-Windows" {
			t.Error("INCORRECT LABEL", labels[i])
			break
		}
	}
}

func TestBasicLabelLLVM(t *testing.T) {
  lbModel := LBModel{Classifier: &LBClassifier{Ctx: ctx, Client: client}}
  lbModel.Learn([]string{"A-allocators", "A-LLVM"})
  text := "Assume at least LLVM 3.9 in rustllvm and rustc_llvm"
  Case(t, "A-LLVM", text, lbModel)
}

func TestAdvancedLabelAllocator(t *testing.T) {
  //https://github.com/rust-lang/rust/issues/42025
  lbModel := LBModel{Classifier: &LBClassifier{Ctx: ctx, Client: client}}
  lbModel.Learn([]string{"A-allocators", "Os-Linux"})

  text := `https://github.com/rust-lang/rust/blob/master/src/liballoc_system/lib.rs#L231-L241

When reallocating, it calculates the new aligned pointer which may have a different offset than the old allocation, which means the pointer might not point to the same point in the allocation anymore, which is bad. Should probably be changed to just do a completely new allocation and manually copy over the data.

Note that as long as there is no stable way to specify an alignment higher than 8 or 16 bytes (for x86 and x86_64 respectively) this bug cannot be encountered in stable Rust, however it is still important to fix.

Example that reproduces the issue:

#![feature(attr_literals, repr_align)]

#[repr(align(256))]
struct Foo(usize);

fn main() {
    let mut foo = vec![Foo(273)];
    for i in 0..0x1000 {
        foo.reserve_exact(i);
        assert!(foo[0].0 == 273);
        assert!(foo.as_ptr() as usize & 0xff == 0);
        foo.shrink_to_fit();
        assert!(foo[0].0 == 273);
        assert!(foo.as_ptr() as usize & 0xff == 0);
    }
}`
  Case(t, "A-allocators", text, lbModel)
}


func TestAdvancedLabelWindows(t *testing.T) {
	lbModel := LBModel{Classifier: &LBClassifier{Ctx: ctx, Client: client}}
	lbModel.Learn([]string{"Os-Windows", "Os-Linux"})

	text := `Do you want to request a feature or report a bug?
Bug.

What is the current behavior?
I tried installing yarn using the installation script from https://yarnpkg.com/en/docs/install#alternatives-tab (I'm on Windows, but I don't have admin rights, so I can't use the Windows installer). After downloading the script it fails:

$ curl -o- -L https://yarnpkg.com/install.sh | bash
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  6742    0  6742    0     0   1694      0 --:--:--  0:00:03 --:--:--  1701
Installing Yarn!
/c/Users/myuser/.yarn/bin/yarn
stdin is not a tty
When I try to run the install.sh file directly, I get the same error:

$ "./install.sh"
Installing Yarn!
/c/Users/myuser/.yarn/bin/yarn
stdout is not a tty
If the current behavior is a bug, please provide the steps to reproduce.

Download the install.sh file from https://yarnpkg.com/en/docs/install#alternatives-tab
Run from Git bash
What is the expected behavior?
It should install successfully.

Please mention your node.js, yarn and operating system version.
node.js v6.3.0
yarn v0.19.1
Windows 7 x64
git version 2.11.0.windows.3`

	labels, _ := lbModel.Predict(conf.ExpandedIssue{Issue: conf.CRIssue{github.Issue{Title: github.String(text)}, []int{}, []conf.CRPullRequest{}, github.Bool(false)}})
	if len(labels) == 0 {
		t.Error("INCORRECT LABEL. 0 LABELS RETURNED", "EXPECTING Os-Windows")
	}
	for i := 0; i < len(labels); i++ {
		if labels[i] != "Os-Windows" {
			t.Error("INCORRECT LABEL", labels[i])
			break
		}
	}
}
