# marlog
Simple "Stamp" based logging system in Go.

## Docs

[@godoc.org](https://godoc.org/github.com/andmar/marlog "marlog @ godoc.org")

## Usage

Simple usage example:

```go
log := marlog.MarLog
log.Prefix = "TEST"
log.Flags = marlog.FlagLdate | marlog.FlagLtime | marlog.FlagLlongfile

err := log.SetOutputHandle("STDOUT", os.Stout) // NOTE: This would not be needed because os.Stdout is added as a default handle (*STDOUT) on init(). Also a Stamp (*STDOUT) with that handle is also added. These are available to use normally and also via the LogQ method
err = log.SetOutputHandle("STDERR", os.Stderr)
err = log.SetOutputHandle("DISCARD", ioutil.Discard)
if err != nil {
	fmt.Println("Error", err)
}

logFile, err := os.OpenFile("test.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
if err != nil {
	fmt.Println("Could not create the test.log file")
} else {
	log.SetOutputHandle("AFILE", logFile)
}

err = log.SetStamp("DEBUG", "STDOUT", "AFILE")
err = log.SetStamp("ERROR", "AFILE", "*STDOUT") // NOTE: "*STDOUT" is a preset Handle to os.Stdout

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
if err != nil {
	fmt.Println("Error", err)
}

err = log.LogO("ERROR", "This is the sixth logged message", marlog.OptionFatal)
if err != nil {
	fmt.Println("Error:", err)
}

fmt.Println("The tenth logged message should not be...")
```

Using version at commit `51ac76135b57dbf6f79596c4faa1ba95a9f2f6ed`
