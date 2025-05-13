package cchttpheaders

type HTTPHeaderKey int

const (
	NotSet HTTPHeaderKey = iota
	Accept
	RequestTracker
	SessionTracker
	ContentType
	Originator
	RemoteAddress
	Authorization
	RequestID
	XMethod
	XOperator
)

var httpHeaderNames = map[HTTPHeaderKey]string{
	NotSet:         "",
	Accept:         "Accept",
	SessionTracker: "Session-Tracker",
	ContentType:    "Content-Type",
	Originator:     "Originator",
	RemoteAddress:  "Remote-Address",
	Authorization:  "Authorization",
	RequestID:      "Request-Id",
	XMethod:        "X-Method",
	XOperator:      "X-Operator",
}

func (t HTTPHeaderKey) Name() string {
	return httpHeaderNames[t]
}
