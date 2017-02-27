# Heupr&trade; Models

## Description

All **back-end logic** powering the Heupr&trade; project is stored in
this directory. Each model is separated into a unique directory which is titled
the name of the given model. In addition to necessary code assets, these
directories will contain the necessary READMEs, documentation, and resource
references regarding the project.  

**Global assets** that are common across all / most models are stored in
separate files / directories  within the `models/` directory. Several of the
current global assets include:
- tossing graph logic
- cross-fold validation methods
- confusion matrix generators
- internal "simple" issue structs
- "classifier" interface
- broad "parent" model struct
How these are handled on a package management level will be gradually
determined (e.g. directly within `models/` or in specific subdirectories).  
