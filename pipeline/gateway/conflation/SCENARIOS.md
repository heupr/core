# Scenarios

## Description

This is a catalog of all scenario packets that are available for use within
the core data filter. Note that the status of each is also available
in the description.  

## Breakdown

- **ScenarioAND** - pending review
  - "Meta" filter
  - Provides "AND" logic between given scenarios
  - Alternative to built-in "OR" logic in Conflator
- **Scenario1** - not started
  - "Basic" issues
  - Issue objects without any additional criteria
  - Only issues (pull requests excluded from filtering)
- **Scenario2** - incomplete
  - Conversation issues
  - Issues with comment activity
- **Scenario3** - incomplete
  - Closing pull requests
  - Pull requests that officially close a raised issue
- **Scenario4** - pending review
  - "Naked" pull requests
  - Only pull requests without an associated issues
- **Scenario5** - pending review
  - Issue body length
  - User-determined issue body text word count
- **Scenario6** - pending review
  - Number of assignees
  - User-defined number of assignees per issue
- **Scenario7** - pending review
  - Conflating "naked" pull requests
  - Fills the reference fields for conflation
  - Note: this is necessary in the Bhattacharya model
