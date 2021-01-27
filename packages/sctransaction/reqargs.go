package sctransaction

import (
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/registry"
	"io"
)

// RequestArguments encodes request parameters taking into account hashes of data blobs
type RequestArguments dict.Dict

// NewRequestArguments constructor
func NewRequestArguments() RequestArguments {
	return RequestArguments(dict.New())
}

// AddArgument add new ordinary argument. Encodes the key as "normal"
func (a RequestArguments) AddArgument(name string, data []byte) {
	a[kv.Key("-"+name)] = data
}

// AddArgumentAsBlobHash adds argument with the data hash instead of data itself.
// Encodes key as "blob reference"
func (a RequestArguments) AddArgumentAsBlobHash(name string, data []byte) {
	h := hashing.HashData(data)
	a[kv.Key("*"+name)] = h[:]
}

// HasBlogRef return if request arguments contain at least one blog reference
func (a RequestArguments) HasBlogRef() bool {
	var ret bool
	(dict.Dict(a)).ForEach(func(key kv.Key, _ []byte) bool {
		ret = []byte(key)[0] == '*'
		if ret {
			return false
		}
		return true
	})
	return ret
}

func (a RequestArguments) Write(w io.Writer) error {
	return (dict.Dict(a)).Write(w)
}

func (a RequestArguments) Read(r io.Reader) error {
	return (dict.Dict(a)).Read(r)
}

// DecodeRequestArguments decodes RequestArguments. For each blog reference encoded it
// looks for the data by hash into the registry and replaces dict entry with the data
// It returns ok flag == false if at least one blob hash don't have data in the registry
func (a RequestArguments) DecodeRequestArguments(reg *registry.Impl) (dict.Dict, bool, error) {
	ret := dict.New()
	ok := true
	var err error
	var data []byte
	var h hashing.HashValue
	reqArgsDict := dict.Dict(a)
	reqArgsDict.ForEach(func(key kv.Key, value []byte) bool {
		d := []byte(key)
		if d[0] == '*' {
			h, err = hashing.HashValueFromBytes(value)
			if err != nil {
				ok = false
				return false
			}
			data, ok, err = reg.GetBlob(h)
			if err != nil {
				ok = false
				return false
			}
			if !ok {
				ok = false
				return false
			}
			ret.Set(kv.Key(d[1:]), data)
		} else {
			ret.Set(kv.Key(d[1:]), value)
		}
		return true
	})
	if err != nil || !ok {
		ret = nil
	}
	return ret, ok, err
}