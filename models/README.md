# CoralReefCI&trade; Models

## Description

All back-end logic powering the CoralReefCI&trade; project is stored in
this directory. Each model is separated into a unique directory which is titled
the name of the given model.  

Global assets that are common across all / most models are stored in separate
files / directories  within the `models/` directory.  

## Catalog

#### Bhattacharya

This is the initial model built and focuses on utilizing Naive Bayes
classifiers, bug tossing graphs, and other tools, as a means of identifying
the appropriate developer to assign an issue to. Bhattacharya will primarily
stand as a benchmark and possibly a testing ground for future models.  

- **Status**: pre-alpha  
  - the model is currently incomplete and missing various functionality
  - several code assets here need to be refractored into global assets
- **Goals**: benchmarking and testing
  - Bhattacharya will primarily serve as a laboratory environment once
  the model is replaced in production
- **Assumptions**:
  - issue and pull request body text can be used to predict future issue
  assignment
