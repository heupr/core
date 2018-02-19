package ingestor

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
