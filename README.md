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

err = log.LogS("DEBUG", "This is the first logged message")
err = log.LogO("DEBUG", "This is the second logged message", marlog.OptionNone)
err = log.LogC(false, "DEBUG", "The third logged message should not be...")
err = log.LogC(true, "DEBUG", "This is the fourth logged message")
err = log.Log(true, "DEBUG", "This is the fifth logged message", marlog.OptionNone)

log.DeactivateStamps("DEBUG")

err = log.LogS("DEBUG", "The sixth logged message should not be...")
err = log.LogS("ERROR", "The seventh logged message should only go into the STDOUT handle")

log.ActivateStamps("DEBUG")

err = log.LogS("DEBUG", "This is the eigth logged message")
err = log.LogS("ERROR", "This is the nineth logged message")

err = log.LogO("ERROR", "This is the sixth logged message", marlog.OptionFatal)

if err != nil {
	fmt.Println("Error:", err)
}

fmt.Println("The tenth logged message should not be...")
```

Using version at commit `92f6ad82b5b0b7931209539e74cd0981d08074c7`
