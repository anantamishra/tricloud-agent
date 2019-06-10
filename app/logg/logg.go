package logg

import "log"

var logging bool = true

func Log(v ...interface{}) {
	if !logging {
		return
	}

	//log.Println(v...)
}

func Debug(v ...interface{}) {

	log.Println(v...)
}
