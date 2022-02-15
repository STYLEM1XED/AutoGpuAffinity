package main

import (
	"flag"
)

func init() {
	flag.IntVar(&totaltrials, "totaltrials", -1, "Enter how many trials you would like to test for each CPU")
	flag.IntVar(&trialtime, "trialtime", -1, "Enter how many seconds you want each trial to last")

	flag.Parse()
	if totaltrials != -1 || trialtime != -1 {
		cliMode = true
		if totaltrials == -1 {
			totaltrials = 3 // default
		}
		if trialtime == -1 {
			trialtime = 30 // default
		}
	}
}
