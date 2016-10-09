# Bhattacharya Model

## Description

This model takes its name from the scholarly article by
**Bhattacharya et al.** on which the bulk of the initial model structure is
built. The primary features of this particular model include a Naive Bayes
classifier, tossing graphs, and incremental learning among others.

## Status

Currently, the **Bhattacharya** model is in an operational state. However,
several things must be noted about the various aspects of the model - these
aspects have had notable impacts on the overall success rate of the
prediction results.

1. The model can produce a fully prediction set.
2. `shuffle.go` appears to be broken as it severely negatively
impacted results.
  - This method dropped accuracy down to below 1% when included.
  - There is likely a bug in the logic of the code.
3. `stemmer.go` also had a negative impact although it was not
as notable as `shuffle.go`.
  - Including this method dropped accuracy down to ~35%.
  - Due to this being a 3rd party library, the weakness is
  likely there.
  - It could be due to how the library handles punctuation.
    - E.g. at the end of a sentence.
4. `fold.go` has not been included into the model.
  - This provided incremental learning to the model.

The "straight up" version of the model (without any of the
  additional functions) achieved an accuracy of ~62%. The
  bottom line is that the model can operate but it is likely
  not providing its most optimal results.  

## Forward

At some point in the indefinite future, production will circle
around to complete the model in its entirety for posterity's
sake. Ideally the project's output will be used as a solid
benchmark for other models. Additionally, several aspects of
this model's code will be refactored out to become global
assets. These include, but are not limited to:

- `confuse.go` - for result evaluation
- `shuffle.go` - to randomize training inputs
