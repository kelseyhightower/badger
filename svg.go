package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const (
	green  = "#34eb62"
	red    = "#e31b14"
	yellow = "#ebd634"
)

const (
	statusFailure = "failure"
	statusSuccess = "success"
	statusUnknown = "unknown"
	statusWorking = "working"
)

const svgTmpl = `
<svg xmlns="http://www.w3.org/2000/svg" width="140" height="20">
  <linearGradient id="a" x2="0" y2="100%%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <rect rx="3" width="140" height="20" fill="#555"/>
  <rect rx="3" x="80" width="60" height="20" fill="%s"/>
  <path fill="%s" d="M80 0h4v20h-4z"/>
  <rect rx="3" width="140" height="20" fill="url(#a)"/>
  <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
    <text x="40" y="15" fill="#010101" fill-opacity=".3">cloud build</text>
    <text x="40" y="14">cloud build</text>
    <text x="110" y="15" fill="#010101" fill-opacity=".3">%s</text>
    <text x="110" y="14">%s</text>
  </g>
</svg>
`

var (
	svgFailure = []byte(fmt.Sprintf(svgTmpl, red, red, statusFailure, statusFailure))
	svgSuccess = []byte(fmt.Sprintf(svgTmpl, green, green, statusSuccess, statusSuccess))
	svgUnknown = []byte(fmt.Sprintf(svgTmpl, red, red, statusUnknown, statusUnknown))
	svgWorking = []byte(fmt.Sprintf(svgTmpl, yellow, yellow, statusWorking, statusWorking))
)

var (
	svgFailureHash = hash(svgFailure)
	svgSuccessHash = hash(svgSuccess)
	svgUnknownHash = hash(svgUnknown)
	svgWorkingHash = hash(svgWorking)
)

func hash(data []byte) string {
	s := md5.Sum(data)
	return hex.EncodeToString(s[:])
}
