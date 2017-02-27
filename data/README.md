# Heupr&trade; Data Storage

## Description

This repository holds the various files generated throughout the course of
running the program. See below for details on the three categories out output.

### Backtests

**Backtests** include the actual output results of running the various model
backtests. This will include information such as the model name, included
scenarios, training / testing data counts, and others.  

### Caches

**Caches** holds the output from data pulldowns from **GitHub** so as to avoid
overloading calls to the API. Information here will be structured according to
how it is organized by the **go-github** third-party API interface.  

### Logs

**Logs** contains the raw output from the various functions of the model. This
is used to fine tune / identify areas for improving the model. As an example,
data collected here will include the name of the logger (which should provide
information on what is being evaluated), run time, output, and others.  
