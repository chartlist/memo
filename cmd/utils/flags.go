/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-10-11
 */
package utils

import (
	"strings"

	"github.com/urfave/cli"
)

var (
	ConfigFlag = cli.StringFlag{
		Name:  "config",
		Usage: "Genesis block config `<file>`. If doesn't specifies, use main net config as default.",
	}

	FromFlag = cli.Uint64Flag{
		Name:  "from",
		Value: 0,
		Usage: "",
	}

	ToFlag = cli.Uint64Flag{
		Name:  "to",
		Value: 0,
		Usage: "",
	}
)

func GetFlagName(flag cli.Flag) string {
	name := flag.GetName()
	if name == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(name, ",")[0])
}
