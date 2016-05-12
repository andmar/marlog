# TODO
 * Allow for logging to multiple Stamps on the same call
 * Print information about who (e.g. which function) Logged the line: Help at http://stackoverflow.com/questions/35212985/is-it-possible-get-information-about-caller-function-in-golang#35213181, by using the Log package what is outputted is the Log systems file, which is irrelevant. See more information on the Log package code? They do this.
 * Consider making the Logging functions package based instead of attached to the main object so the "calling code's line" is smaller
 * If provided Stamp/Handler does not exist? Default to log to stdout?
 * Remove "Log" from the log functions to reduce verbosity and "log.Log..." usage weirdness
