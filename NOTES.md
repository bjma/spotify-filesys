# Production Notes
* Should make printing tracks faster; we should note some overheads
    * Fetching trees itself will take a lot of time. 
    We *could* optimize the function in the `api` package to solve this
    * There is overhead in initializing the tracks at tree construction as 
    well, so we'd have to weigh the advantages of fetching on demand
    or fetching on construction
