# Heupr Tests

## Description

This directory contains the various **testing structures** necessary for the
Heupr project. These will include several types such as model
backtesting, stress testing infrastructure, handling failure scenarios, and
others (as they are developed).  

All of the results from these tests are being output as files into the data/
directory (not contained within version control). Note that this directory
focuses on a cmd/ directory that houses all of the subdirectories clustered as
individual `main` packages holding the necessary executables to fire the
desired test packets. Each directory is named after its corresponding model /
project aspect (e.g. *bhattacharya*).  

Think of this directory as containing "model-level" or "project-level" unit
tests. There may also eventually be tests available at the "asset-level" of the
project.  
