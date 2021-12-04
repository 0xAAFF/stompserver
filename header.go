package stompserver

import (
	"strconv"
)

// STOMP header names. Some of the header
// names have commands with the same name
// (eg Ack, Message, Receipt). Commands use
// an upper-case naming convention, header
// names use pascal-case naming convention.
const (
	AcceptVersion = "accept-version"
	Ack           = "ack"
	ContentLength = "content-length"
	ContentType   = "content-type"
	Destination   = "destination"
	HeartBeat     = "heart-beat"
	Host          = "host"
	Id            = "id"
	Login         = "login"
	Message       = "message"
	MessageId     = "message-id"
	Passcode      = "passcode"
	Receipt       = "receipt"
	ReceiptId     = "receipt-id"
	Server        = "server"
	Session       = "session"
	Subscription  = "subscription"
	Transaction   = "transaction"
	Version       = "version"
)

type headerKey struct {
	AcceptVersion string
	Ack           string
	ContentLength string
	ContentType   string
	Destination   string
	HeartBeat     string
	Host          string
	Id            string
	Login         string
	Message       string
	MessageId     string
	Passcode      string
	Receipt       string
	ReceiptId     string
	Server        string
	Session       string
	Subscription  string
	Transaction   string
	Version       string
}

var StompHeaders = &headerKey{
	AcceptVersion: "accept-version",
	Ack:           "ack",
	ContentLength: "content-length",
	ContentType:   "content-type",
	Destination:   "destination",
	HeartBeat:     "heart-beat",
	Host:          "host",
	Id:            "id",
	Login:         "login",
	Message:       "message",
	MessageId:     "message-id",
	Passcode:      "passcode",
	Receipt:       "receipt",
	ReceiptId:     "receipt-id",
	Server:        "server",
	Session:       "session",
	Subscription:  "subscription",
	Transaction:   "transaction",
	Version:       "version",
}

// A Header represents the header part of a STOMP frame.
// The header in a STOMP frame consists of a list of header entries.
// Each header entry is a key/value pair of strings.
//
// Normally a STOMP header only has one header entry for a given key, but
// the STOMP standard does allow for multiple header entries with the same
// key. In this case, the first header entry contains the value, and any
// subsequent header entries with the same key are ignored.
//
// Example header containing 6 header entries. Note that the second
// header entry with the key "comment" would be ignored.
//
//	login:scott
//	passcode:tiger
//	host:stompserver
//	accept-version:1.0,1.1,1.2
//	comment:some comment
//	comment:another comment
//
type Header struct {
	slice []string
}

// NewHeader creates a new Header and populates it with header entries.
// This function expects an even number of strings as parameters. The
// even numbered indices are keys and the odd indices are values. See
// the example for more information.
func NewHeader(headerEntries ...string) *Header {
	h := &Header{}
	h.slice = append(h.slice, headerEntries...)
	if len(h.slice)%2 != 0 {
		h.slice = append(h.slice, "")
	}
	return h
}

// Add adds the key, value pair to the header.
func (h *Header) Add(key, value string) {
	h.slice = append(h.slice, key, value)
}

// AddHeader adds all of the key value pairs in header to h.
func (h *Header) AddHeader(header *Header) {
	if header != nil {
		for i := 0; i < header.Len(); i++ {
			key, value := header.GetAt(i)
			h.Add(key, value)
		}
	}
}

// AddFromArray adds all of the key value pairs in header to h.
func (h *Header) AddFromArray(array []string) {
	if array != nil && len(array)%2 == 0 {
		for i := 0; i < len(array); i += 2 {
			h.slice = append(h.slice, array[i], array[i+1])
		}
	}
}

// Set replaces the value of any existing header entry with the specified key.
// If there is no existing header entry with the specified key, a new
// header entry is added.
func (h *Header) Set(key, value string) {
	if i, ok := h.index(key); ok {
		h.slice[i+1] = value
	} else {
		h.slice = append(h.slice, key, value)
	}
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns "".
func (h *Header) Get(key string) string {
	value, ok := h.Contains(key)
	if ok {
		return value
	} else {
		return ""
	}

}

// GetAll returns all of the values associated with a given key.
// Normally there is only one header entry per key, but it is permitted
// to have multiple entries according to the STOMP standard.
func (h *Header) GetAll(key string) []string {
	var values []string
	for i := 0; i < len(h.slice); i += 2 {
		if h.slice[i] == key {
			values = append(values, h.slice[i+1])
		}
	}
	return values
}

// Returns the header name and value at the specified index in
// the collection. The index should be in the range 0 <= index < Len(),
// a panic will occur if it is outside this range.
func (h *Header) GetAt(index int) (key, value string) {
	index *= 2
	return h.slice[index], h.slice[index+1]
}

// Contains gets the first value associated with the given key,
// and also returns a bool indicating whether the header entry
// exists.
//
// If there are no values associated with the key, Get returns ""
// for the value, and ok is false.
func (h *Header) Contains(key string) (value string, ok bool) {
	var i int
	if i, ok = h.index(key); ok {
		value = h.slice[i+1]
	}
	return
}

// Contains Key
func (h *Header) ContainsKey(key string) bool {
	for i := 0; i < len(h.slice); i = i + 2 {
		if h.slice[i] == key {
			return true
		}
	}
	return false
}

// Del deletes all header entries with the specified key.
func (h *Header) Del(key string) {
	for i, ok := h.index(key); ok; i, ok = h.index(key) {
		h.slice = append(h.slice[:i], h.slice[i+2:]...)
	}
}

// Len returns the number of header entries in the header.
func (h *Header) Len() int {
	return len(h.slice) / 2
}

// Clone returns a deep copy of a Header.
func (h *Header) Clone() *Header {
	hc := &Header{slice: make([]string, len(h.slice))}
	copy(hc.slice, h.slice)
	return hc
}

// ContentLength returns the value of the "content-length" header entry.
// If the "content-length" header is missing, then ok is false. If the
// "content-length" entry is present but is not a valid non-negative integer
// then err is non-nil.
func (h *Header) ContentLength() (value int, ok bool, err error) {
	text, ok := h.Contains(ContentLength)
	if !ok {
		return 0, false, nil
	}

	n, err := strconv.ParseUint(text, 10, 32)
	if err != nil {
		return 0, true, err
	}

	value = int(n)
	ok = true
	return value, ok, nil
}

// Returns the index of a header key in Headers, and a bool to indicate
// whether it was found or not.
func (h *Header) index(key string) (int, bool) {
	for i := 0; i < len(h.slice); i += 2 {
		if h.slice[i] == key {
			return i, true
		}
	}
	return -1, false
}
