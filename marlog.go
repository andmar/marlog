package marlog

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"math/rand"

	"gitlab.com/voipit/thalesecscore/ecscoreapi/utils"
)

const (
	defaultStampName        = "*STDOUT"
	defaultOutputHandleName = "*STDOUT"
)

const (
	// FlagLdate Refer to the Log package documentation
	FlagLdate = log.Ldate
	// FlagLtime Refer to the Log package documentation
	FlagLtime = log.Ltime
	// FlagLmicroseconds Refer to the Log package documentation
	FlagLmicroseconds = log.Lmicroseconds
	// FlagLlongfile Refer to the Log package documentation
	FlagLlongfile = log.Llongfile
	// FlagLshortfile Refer to the Log package documentation
	FlagLshortfile = log.Lshortfile
	// FlagLUTC Refer to the Log package documentation
	FlagLUTC = log.LUTC
	// FlagLstdFlags Refer to the Log package documentation
	FlagLstdFlags = log.Ldate | log.Ltime
)

const (
	// OptionNone Use this option or 0 as the options value when calling the Log... methods and you don't want to pass any options
	OptionNone = 1 << iota
	// OptionFatal This option makes the Log... methods call os.Exit(-1) after printing the log message
	OptionFatal
)

// MarLog Variable with precreated MarLogger
var MarLog *MarLogger

func init() {
	MarLog = new(MarLogger)
	MarLog.Prefix = ""
	MarLog.Flags = FlagLstdFlags
	MarLog.Active = true
	MarLog.stamps = make(map[string]*stamp)
	MarLog.outputHandles = make(map[string]*outputHandle)

	MarLog.lockCond = sync.NewCond(&sync.Mutex{})
	MarLog.aquireMU = &sync.Mutex{}
	MarLog.releaseMU = &sync.Mutex{}

	MarLog.SetOutputHandle(defaultOutputHandleName, os.Stdout)
	MarLog.SetStamp(defaultStampName, defaultOutputHandleName)
}

// MarLogger The MarLogger type
type MarLogger struct {
	Prefix           string
	Flags            int
	Active           bool
	stamps           map[string]*stamp
	outputHandles    map[string]*outputHandle
	lockingContextID string
	lockCond         *sync.Cond
	aquireMU         *sync.Mutex
	releaseMU        *sync.Mutex
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

// LogSQ (Quick) Print a log line with a specific message to stdout (aliast to fmt.Println)
func (logger *MarLogger) LogSQ(message string) error {
	return logger.Log(context.TODO(), true, defaultStampName, message, OptionNone)
}

// LogSS (Simple) Print a log line with a specific message to the output handles of a specific stamp
func (logger *MarLogger) LogSS(stampName string, message string) error {
	return logger.Log(context.TODO(), true, stampName, message, OptionNone)
}

// LogSO (Option) Print a log line with a specific message to the output handles of a specific stamp with options
func (logger *MarLogger) LogSO(stampName string, message string, options int) error {
	return logger.Log(context.TODO(), true, stampName, message, options)
}

// LogSC (Condition) Print a log line with a specific message to the output handles of a specific stamp, if the condition is true
func (logger *MarLogger) LogSC(condition bool, stampName string, message string) error {
	return logger.Log(context.TODO(), condition, stampName, message, OptionNone)
}

// LogQ (Quick) Print a log line with a specific message to stdout (aliast to fmt.Println). Context related features available.
func (logger *MarLogger) LogQ(ctx context.Context, message string) error {
	return logger.Log(ctx, true, defaultStampName, message, OptionNone)
}

// LogS (Simple) Print a log line with a specific message to the output handles of a specific stamp. Context related features available.
func (logger *MarLogger) LogS(ctx context.Context, stampName string, message string) error {
	return logger.Log(ctx, true, stampName, message, OptionNone)
}

// LogO (Option) Print a log line with a specific message to the output handles of a specific stamp with options. Context related features available.
func (logger *MarLogger) LogO(ctx context.Context, stampName string, message string, options int) error {
	return logger.Log(ctx, true, stampName, message, options)
}

// LogC (Condition) Print a log line with a specific message to the output handles of a specific stamp, if the condition is true. Context related features available.
func (logger *MarLogger) LogC(ctx context.Context, condition bool, stampName string, message string) error {
	return logger.Log(ctx, condition, stampName, message, OptionNone)
}

// Log Print a log line with a specific message to the output handles of a specific stamp with options, if the condition is true
func (logger *MarLogger) Log(ctx context.Context, condition bool, stampName string, message string, options int) error {

	if !logger.Active {
		return nil
	}

	contextid := ""
	switch value := ctx.Value("ID").(type) {
	case string:
		contextid = value
	}

	if contextid != logger.lockingContextID && logger.lockingContextID != "" {
		logger.LockContext(utils.CreateContextID())
		defer logger.UnlockContext()
	}

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

			prefix := ""
			if logger.Prefix != "" {
				prefix = logger.Prefix + " "
			}

			log := log.New(outputHandle.handle, prefix, logger.Flags&^FlagLshortfile&^FlagLlongfile)

			newLine := "\n"
			if options&OptionFatal != 0 {
				newLine = ""
			}

			// NOTE: Get Calling function information, refer to https://github.com/golang/go/blob/master/src/log/log.go for more information on this
			_, file, line, ok := runtime.Caller(2)
			if !ok {
				return fmt.Errorf("Error getting Caller information")
			}
			fileSplitSlice := strings.Split(file, "/")
			fileShort := fileSplitSlice[len(fileSplitSlice)-1]

			if contextid != "" {
				message = "(Context:" + contextid + ") " + message
			}

			if stamp.MessagePrefix != "" {
				if logger.Flags&FlagLlongfile != 0 {
					log.Printf("%s (%v) %s: %s%s", file, line, stamp.MessagePrefix, message, newLine)
				} else if logger.Flags&FlagLshortfile != 0 {
					log.Printf("%s (%v) %s: %s%s", fileShort, line, stamp.MessagePrefix, message, newLine)
				} else {
					log.Printf("%s: %s%s", stamp.MessagePrefix, message, newLine)
				}
			} else {
				log.Printf("%s (%v): %s%s", file, line, message, newLine)
			}

		}

		if options&OptionFatal != 0 { // NOTE: This is kinda the same as using Golang's log.Fatalf
			os.Exit(-1)
		}

	}

	return nil

}

// SetStamp Try to setup a new Stamp with the specific name using the specified output handles
func (logger *MarLogger) SetStamp(stampName string, outputHandleKeys ...string) error {

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
	newStamp.MessagePrefix = stampName
	newStamp.HandleKeys = handleKeys
	logger.stamps[stampName] = newStamp

	return nil
}

// SetOutputHandle Try to setup a new OutputHandle with the specified name and io handle
func (logger *MarLogger) SetOutputHandle(handleName string, handle io.Writer) error {

	if _, found := logger.outputHandles[handleName]; found == true {
		return fmt.Errorf("Output Handle named \"%s\" already exists.", handleName)
	}

	newOutputHandle := new(outputHandle)
	newOutputHandle.Name = handleName
	newOutputHandle.handle = handle
	logger.outputHandles[handleName] = newOutputHandle

	return nil
}

// AddOuputHandles Try to add output handles to a Stamp
func (logger *MarLogger) AddOuputHandles(stampName string, outputHandleKeys ...string) error {

	stamp, found := logger.stamps[stampName]
	if found == false {
		return fmt.Errorf("Stamp named \"%s\" does not exist.", stampName)
	}

	for _, currentHandleKey := range outputHandleKeys {
		stamp.HandleKeys = append(stamp.HandleKeys, currentHandleKey)
	}

	return nil
}

// ActivateStamps Sets the Stamps with names passed as arguments as active
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

// DeactivateStamps Sets the Stamps with names passed as arguments as not active
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

// LockContext Locks the LockCond associated with the logger to a context with a ID key equal to contextID, while waiting to have exclusive access to main log function.
func (logger *MarLogger) LockContext(contextID string) {

	tryAquire := func(id string) bool {

		logger.aquireMU.Lock()
		defer logger.aquireMU.Unlock()

		if logger.lockingContextID == "" {
			logger.lockingContextID = id
			return true
		}

		return false

	}

	logger.lockCond.L.Lock()
	defer logger.lockCond.L.Unlock()
	for !tryAquire(contextID) {
		fmt.Println("Blocked. Trying to aquire lock:", contextID, "Current Holder:", logger.lockingContextID)
		logger.lockCond.Wait()
	}

	fmt.Println("Proceding...", logger.lockingContextID)

}

// UnlockContext Unlocks all goroutines waiting on the LockConf associated with the logger.
func (logger *MarLogger) UnlockContext() {

	logger.releaseMU.Lock()
	defer logger.releaseMU.Unlock()

	fmt.Println("Releasing Lock. Holder:", logger.lockingContextID)

	logger.lockingContextID = ""
	logger.lockCond.Broadcast()

}

// CreateContextID Creates a string to be used as Context ID value for logging contexts where no ID is provided via the context object
func createContextID() string {

	return fmt.Sprint(time.Now().Unix(), "-", randStringRunes(5))

}

// randStringRunes Returns a string of size n with random characters
func randStringRunes(n int) string {

	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)

}
