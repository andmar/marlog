# marlog
Simple "Stamp" based logging system in Go

## Usage

Simple usage example:

```go
log := marlog.MarLog
log.Prefix = "Test"

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

err = log.Log("DEBUG", "This is the first logged message", true, false)
err = log.Log("ERROR", "This is the second logged message", true, false)

log.DeactivateStamps("DEBUG")

err = log.Log("DEBUG", "This is the third logged message", true, false)
err = log.Log("ERROR", "This is the fourth logged message", true, false)

log.ActivateStamps("DEBUG")

err = log.Log("DEBUG", "This is the fifth logged message", true, false)
err = log.Log("ERROR", "This is the sixth logged message", true, false)

if err != nil {
	fmt.Println("Error", err)
}
```

Using version at commit `77e7f177ba6cf288560165f94d66b6ba12990ba1`
