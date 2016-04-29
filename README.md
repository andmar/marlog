# marlog
Simple "Stamp" based logging system in Go

## Usage

Simple usage example:

```go
log := marlog.MarLog
log.Prefix = "TEST"

err := log.AddOutputHandle("STDOUT", os.Stdout)
err = log.AddOutputHandle("STDERR", os.Stderr)
err = log.AddOutputHandle("DISCARD", ioutil.Discard)

if err != nil {
	fmt.Println("Error", err)
}

logFile, err := os.OpenFile("test.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
if err != nil {
	fmt.Println("Could not create the test.log file")
} else {
	log.AddOutputHandle("FILEA", logFile)
}

err = log.AddStamp("DEBUG", "STDOUT", "FILEA")
err = log.AddStamp("ERROR", "STDOUT")

if err != nil {
	fmt.Println("Error", err)
}

fmt.Println(marlog.MarLog)

err = log.Log("DEBUG", "This is the first logged message", 0)
err = log.Log("ERROR", "This is the second logged message", 0)

log.DeactivateStamps("DEBUG")

err = log.Log("DEBUG", "This is the third logged message", 0)
err = log.Log("ERROR", "This is the fourth logged message", 0)

log.ActivateStamps("DEBUG")

err = log.Log("DEBUG", "This is the fifth logged message", 0)
err = log.Log("ERROR", "This is the sixth logged message", 0)

if err != nil {
	fmt.Println("Error", err)
}

err = log.Log("DEBUG", "This is the seventh logged message", 0)
err = log.Log("DEBUG", "This is the eighth logged message", marlog.OptionFatal)

if err != nil {
	fmt.Println("Error", err)
}

fmt.Println("The nineth logged message should not be...")
```

Using version at commit `f1125dad46b9fccc0c8ad396eaa56aaacc0de244`
