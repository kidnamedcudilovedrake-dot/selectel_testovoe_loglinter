package b

import (
	"log"
)

func TestCustomSensitive(secret_val string, key_99 string, normalVar string) {
	log.Println(secret_val) // want "log call contains potentially sensitive variable \"secret_val\""
	log.Println(key_99)     // want "log call contains potentially sensitive variable \"key_99\""
	log.Println(normalVar)
}
