package gentzen

import (
	"log"
	"math/rand"
	"time"
)

type config struct {
	oPL,
	oTHM,
	oCOND,
	oML,
	oDR,
	oDEBUG bool
}

var (
	oPL    = false
	oTHM   = false
	oCOND  = true
	oAX    = true
	oML    = false
	oDR    = false
	oDEBUG = false
	oSYMB  = false //false means standard Polish
)

func init() {

	rand.Seed(int64(time.Now().Nanosecond()))
	SetStandardPolish(true)
	//connectives = append(connectivesSL, connectivesPL...)
	logger = log.New(&checkLog, "", 0)
	debug = log.New(&dLog, "gentzen: ", 0)

}

// SetConditional determines whether we use the conditional.
func SetConditional(v bool) {

	oCOND = v

	setupConnectives()

	return
}

// SetPL enables Predicate Logic Processing when v is true.
// By default, PL processing is disabled.
func SetPL(v bool) {

	oPL = v

}

// SetML specifies whether we allows modal logic
func SetML(v bool) {

	oML = v

}

// SetDR specifies whether we allow derived rules
func SetDR(v bool) {

	oDR = v

}

// SetAllowTheorems sets whether appeal to some standard theorems
// is allowed. Default is false.
func SetAllowTheorems(v bool) {
	oTHM = v
}

// SetStrict sets whether inferece rules should be checked strictly.
func SetStrict(v bool) {
	strictCheck = v
}
