package marlog

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	// NOTE: Refer to Log package documentation for the meaning of these, these are just poiting to them.
	// FlagLdate ...
	FlagLdate = log.Ldate
	// FlagLtime ...
	FlagLtime = log.Ltime
	// FlagLmicroseconds ...
	FlagLmicroseconds = log.Lmicroseconds
	// FlagLlongfile ...
	FlagLlongfile = log.Llongfile
	// FlagLshortfile ...
	FlagLshortfile = log.Lshortfile
	// FlagLUTC ...
	FlagLUTC = log.LUTC
	// FlagLstdFlags ...
	FlagLstdFlags = log.Ldate | log.Ltime
)

const (
	// OptionNone Use this option or 0 as the options value when calling the Log... methods and you don't want to pass any options
	OptionNone = 0 << iota
	// OptionFatal This option makes the Log... methods call os.Exit(-1) after printing the log message
	OptionFatal
	// OptionNewLine ...
	//OptionNewLine
)

// MarLog Variable with precreated MarLogger
var MarLog *MarLogger

// MarLogger The MarLogger type
type MarLogger struct {
	Prefix        string
	Flags         int
	stamps        map[string]*stamp
	outputHandles map[string]*outputHandle
}

type stamp struct {
	Name          string // NOTE: Should be the same as the key in the stamps map
	Active        bool
	MessagePrefix string
	HandleKeys    []string // NOTE: Optionally filter the outputHandlers to use. If this is empty the system will use all output handles to output messages with this Stamp
}

type outputHandle struct {
	Name   string // NOTE: Should be the same as the key in the stamps map
	handle io.Writer
}

// LogS Print a log line with a specific message to the output handles of a specific stamp
func (logger *MarLogger) LogS(stampName string, message string) error {
	return logger.Log(true, stampName, message, OptionNone)
}

// LogO Print a log line with a specific message to the output handles of a specific stamp with options
func (logger *MarLogger) LogO(stampName string, message string, options int) error {
	return logger.Log(true, stampName, message, options)
}

// LogC Print a log line with a specific message to the output handles of a specific stamp, if the condition is true
func (logger *MarLogger) LogC(condition bool, stampName string, message string) error {
	return logger.Log(condition, stampName, message, OptionNone)
}

// Log Print a log line with a specific message to the output handles of a specific stamp with options, if the condition is true
func (logger *MarLogger) Log(condition bool, stampName string, message string, options int) error {

	if condition == true {

		if _, found := logger.stamps[stampName]; found == false {
			return fmt.Errorf("Stamp named \"%s\" does not exist.", stampName)
		}

		stamp := logger.stamps[stampName]

		if stamp.Active == false {
			return fmt.Errorf("Stamp named \"%s\" is not active.", stampName)
		}

		for _, currentHandleKey := range stamp.HandleKeys {

			outputHandle, found := logger.outputHandles[currentHandleKey]
			if found == false {
				return fmt.Errorf("Output Handle named \"%s\" not found.", currentHandleKey)
			}
			log := log.New(outputHandle.handle, stamp.MessagePrefix, logger.Flags)

			if logger.Prefix != "" {
				log.Printf("%s: %s\n", logger.Prefix, message)
			} else {
				log.Printf("%s\n", message)
			}

		}

		if options&OptionFatal != 0 { // NOTE: This is kinda the same as using Golang's log.Fatalf
			os.Exit(-1)
		}

		return nil

	}

	return nil

}

// AddStamp ...
func (logger *MarLogger) AddStamp(stampName string, outputHandleKeys ...string) error {

	if _, found := logger.stamps[stampName]; found == true {
		return fmt.Errorf("Stamp named \"%s\" already exists.", stampName)
	}

	var handleKeys []string
	for _, currentHandleKey := range outputHandleKeys {
		handleKeys = append(handleKeys, currentHandleKey)
	}

	newStamp := new(stamp)
	newStamp.Name = stampName
	newStamp.Active = true
	newStamp.MessagePrefix = ""
	newStamp.HandleKeys = handleKeys
	logger.stamps[stampName] = newStamp

	return nil
}

// AddOutputHandle ...
func (logger *MarLogger) AddOutputHandle(handleName string, handle io.Writer) error {

	if _, found := logger.outputHandles[handleName]; found == true {
		return fmt.Errorf("Output Handle named \"%s\" already exists.", handleName)
	}

	newOutputHandle := new(outputHandle)
	newOutputHandle.Name = handleName
	newOutputHandle.handle = handle
	logger.outputHandles[handleName] = newOutputHandle

	return nil
}

// ActivateStamps ...
func (logger *MarLogger) ActivateStamps(stampNames ...string) error {

	if len(stampNames) == 0 {
		return fmt.Errorf("No Stamp names provided")
	}

	for _, currentStampName := range stampNames {
		if _, found := logger.stamps[currentStampName]; found == false {
			return fmt.Errorf("Nothing done. Stamp named \"%s\" does not exist.", currentStampName)
		}
	}

	for _, currentStampName := range stampNames {
		logger.stamps[currentStampName].Active = true
	}

	return nil
}

// DeactivateStamps ...
func (logger *MarLogger) DeactivateStamps(stampNames ...string) error {

	if len(stampNames) == 0 {
		return fmt.Errorf("No Stamp names provided")
	}

	for _, currentStampName := range stampNames {
		if _, found := logger.stamps[currentStampName]; found == false {
			return fmt.Errorf("Nothing done. Stamp named \"%s\" does not exist.", currentStampName)
		}
	}

	for _, currentStampName := range stampNames {
		logger.stamps[currentStampName].Active = false
	}

	return nil
}

func init() {
	MarLog = new(MarLogger)
	MarLog.Prefix = ""
	MarLog.Flags = FlagLdate | FlagLtime | FlagLshortfile
	MarLog.stamps = make(map[string]*stamp)
	MarLog.outputHandles = make(map[string]*outputHandle)
}
