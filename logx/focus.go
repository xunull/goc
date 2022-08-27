package logx

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

func Focus(format string, v ...interface{}) {

	f := "\n" + strings.Repeat("+", 25) + "\n" +
		fmt.Sprintf(format, v...) +
		"\n" + strings.Repeat("-", 25) + "\n"
	// todo \n
	//log.Debug().Msg(f)
	//log.Print(f)
	fmt.Println(f)

}

func FocusErr(err error, format string, v ...interface{}) {
	f := "\n" + strings.Repeat("+", 25) + "\n" +
		format +
		"\n" + strings.Repeat("-", 25) + "\n"
	log.Debug().Err(err).Msgf(f, v...)
}
