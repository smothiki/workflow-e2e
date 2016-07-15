package util

import (
	"math"
	"os"
	"strconv"
	"strings"
)

// PrependError adds 'Error: ' to an expected error, like the CLI does to error messages.
func PrependError(expected error) string {
	return "Error: " + expected.Error()
}

//GetBp return bp array
func GetBp() []string {
	bp := []string{
		"https://github.com/deis/example-clojure-ring.git",
		"https://github.com/deis/example-go.git",
		"https://github.com/deis/example-java-jetty.git",
		"https://github.com/deis/example-nodejs-express.git",
		"https://github.com/deis/example-perl.git",
		"https://github.com/deis/example-php.git",
		"https://github.com/deis/example-play.git",
		"https://github.com/deis/example-python-django.git",
		"https://github.com/deis/example-python-flask.git",
		"https://github.com/deis/example-ruby-sinatra.git",
		"https://github.com/deis/example-scala.git",
	}
	return bp
}

// Getbuildpack gets unique buildpack for each build number
func Getbuildpack() (string, string) {
	bp := GetBp()[GetBpNum()]
	return GetBpName(bp), bp
}

//GetBpName getbuildpack name
func GetBpName(bp string) string {
	return strings.Split(strings.Split(bp, "/")[4], ".")[0]
}

// GetBpNum gets the bp number for that build number
func GetBpNum() int {
	num, _ := strconv.Atoi(os.Getenv("BUILD_NUMBER"))
	return int(math.Mod(float64(num), float64(len(GetBp()))))
}

// GetNextBp get next build pack
func GetNextBp() (string, string) {
	bpn := GetBpNum()
	arr := GetBp()
	if bpn+1 >= len(arr) {
		bpn = 0
	}
	return GetBpName(arr[bpn]), arr[bpn]
}
