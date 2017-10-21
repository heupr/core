package ingestor

const WelcomeTitle = "Hello, I'm Heupr!"

const WelcomeBody = `
## Welcome to the Heupr Integration!

### Introduction
This space is where you will be able to adjust Heupr's settings on how you would like your new ` + "`" + "Issues" + "`" + ` to be handled.

I am programmed to _automatically_ triage newly raised ` + "`" + "Issues" + "`" + `to the appropriate developer based on their _subject matter expertise_. For example, if you'd like me to triage **all** outstanding open ` + "`" + "Issues" + "`" + ` raised since January 2017, I can do that in a jiffy - just let me know!

Please take a moment to review and edit your config settings. Changes to settings need to be made directly in **THIS** issue body and changes take effect when the issue is closed.

### Configuration settings
This is the "settings" page for your Heupr Integration. Please make you selections inside the double quotes in the format indicated.

 - Specific contributors to not assign (e.g. ContributorBlacklist="Bot1, Bot2")
 ContributorBlacklist=""
 - Automated Triage Start Date (e.g. TriageStartDate="01/01/2017")
 TriageStartDate="10/15/2017"
 - Labels to avoid assigning (e.g. IgnoreLabels="Up For Grabs", "Won't Fix")
 IgnoreLabels=""
- Contact info (Optional)
 Email= ""
 Twitter= ""

### Contact
 If you have any questions please reach out to heuprhq@gmail.com.
`

const ConfirmationMessage = `Thank you!. Would like like me to apply the following settings?
                             ContributorBlacklist: %v
                             TriageStartTime: %v
                             IgnoreLabels: %v
                             Email: %v
                             Twitter: %v`
