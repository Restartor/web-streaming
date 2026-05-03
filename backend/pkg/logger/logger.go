package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init() {
	Log = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()

	// penjelasan : Inisialisasi logger dengan output ke stdout dan menambahkan timestamp

}
