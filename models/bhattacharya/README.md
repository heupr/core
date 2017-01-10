# Bhattacharya Model

## Description

This model takes its name from the scholarly article by
**Bhattacharya et al.** on which the bulk of the initial model structure is
built.  

This is the initial model built and focuses on utilizing Naive Bayes
classifiers, bug tossing graphs, incremental learning, and other tools, as a
means of identifying the appropriate developer to assign an issue to.
Bhattacharya will primarily stand as a benchmark and possibly a testing ground
for future models.  

- **Status**: pre-alpha
  - several code assets here need to be refactored into global assets
- **Goals**: benchmarking and testing
  - Bhattacharya will primarily serve as a laboratory environment once
  the model is replaced in production
  - changes to Bhattacharya will not be stored in the `bhattacharya` directory;
  significant changes outside of the paper will constitute a new model
- **Assumptions**:
  - issue and pull request body text can be used to predict future issue
  assignment
  - TF-IDF is an efficient means of identifying training features
  - Naive Bayes correctly classifies features
  - Lapace smoothing is implemented with a smoothing variable of 1
  - regularization has not been implemented in the model
