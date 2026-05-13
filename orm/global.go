package orm

import "time"

var (
	currentDB        Executor
	currentContextDB ContextExecutor
	timestampLoc     *time.Location = time.UTC
)

func SetDB(exec Executor) {
	currentDB = exec
	if ce, ok := exec.(ContextExecutor); ok {
		currentContextDB = ce
	}
}

func GetDB() Executor {
	return currentDB
}

func GetContextDB() ContextExecutor {
	return currentContextDB
}

func SetLocation(loc *time.Location) {
	timestampLoc = loc
}

func GetLocation() *time.Location {
	return timestampLoc
}
