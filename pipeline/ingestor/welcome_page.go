package ingestor

const WelcomeTitle = "Hello, I'm Heupr!"

const WelcomeBody = `
## Welcome to the Heupr Integration!

### Introduction
This space is where you will be able to adjust Heupr's settings on how you would like your new ` + "`" + "Issues" + "`" + ` to be handled.

I am programmed to _automatically_ triage newly raised ` + "`" + "Issues" + "`" + `to the appropriate developer based on their _subject matter expertise_. For example, if you'd like me to triage **all** outstanding open ` + "`" + "Issues" + "`" + ` raised since January 2017, I can do that in a jiffy - just let me know! I'll also place a "triaged" label on each Issue I assign if you have one provided already on your repository.

Please take a moment to review and edit your config settings. Changes to settings need to be made directly in **THIS** issue body.

### Configuration settings
This is the "settings" page for your Heupr Integration. Please make you selections inside the double quotes in the format indicated.

 - Specific users to not assign [comma seperated] (e.g. IgnoreUsers="R2-D2,C-3PO"):
 IgnoreUsers=""
 - Automated Triage Start Date [RFC822] (e.g. TriageStartTime="01 Jan 17 00:00 PST"):
 TriageStartTime="%v"
 - Labels to avoid assigning [comma seperated] (e.g. "wontfix,up for grabs,for padawans"):
 IgnoreLabels=""
 - Contact info (optional, e.g. "darthvader@empire.gov", "@chos3n_0ne"):
 Email=""
 Twitter=""

### Contact
 If you have any questions please reach out to heuprhq@gmail.com.
`

const ConfirmationMessage = `@%v Thank you! Would like me to apply the following settings?[Yes/No]
                             IgnoreUsers: %v
                             TriageStartTime: %v
                             IgnoreLabels: %v
                             Email: %v
                             Twitter: %v`

const ConfirmationErrMessage = `@%v An error occured parsing your settings: %v
                                IgnoreUsers: %v
                                TriageStartTime: %v
                                IgnoreLabels: %v
                                Email: %v
                                Twitter: %v`

const HoldOnMessage = `@%v Understood! The settings will NOT be applied.`

const AppliedSettingsMessage = `@%v Done! Settings processed:
                                IgnoreUsers: %v
                                TriageStartTime: %v
                                IgnoreLabels: %v
                                Email: %v
                                Twitter: %v`
