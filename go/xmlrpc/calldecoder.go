package xmlrpc

import (
	"bytes"
	"encoding/xml"
	"reflect"
)

type calldecoder struct {
	decoder
	nextToken xml.Token
	result    []byte
}

func (dec *calldecoder) Token() (tok xml.Token, err error) {
	if dec.nextToken != nil {
		tok = dec.nextToken
		dec.nextToken = nil
		return tok, nil
	}
	return dec.decoder.Token()
}

func (dec *calldecoder) waitForStartElement(startElementName string) (tok xml.Token, err error) {
	for {
		if tok, err = dec.Token(); err != nil {
			return nil, err
		}
		if t, ok := tok.(xml.StartElement); ok {
			if t.Name.Local != startElementName {
				return nil, invalidXmlError
			}
			return t, nil
		}
		if _, ok := tok.(xml.EndElement); ok {
			dec.nextToken = tok
			return nil, nil
		}
	}
}

func (dec *calldecoder) waitForStartElements(startElementNames ...string) (tok xml.Token, err error) {
	for i, startElementName := range startElementNames {
		tok, err = dec.waitForStartElement(startElementName)
		if err != nil {
			return nil, err
		}
		if nil == tok {
			if i > 0 {
				return nil, invalidXmlError
			}
			return nil, nil
		}
	}
	return tok, nil
}

func (dec *calldecoder) requireStartElements(startElementNames ...string) (tok xml.Token, err error) {
	tok, err = dec.waitForStartElements(startElementNames...)
	if err != nil {
		return nil, err
	}
	if tok == nil {
		return nil, invalidXmlError
	}
	return tok, err
}

func (dec *calldecoder) waitForEndElement(endElementName string) (err error) {
	var tok xml.Token
	for {
		if tok, err = dec.Token(); err != nil {
			return err
		}
		if _, ok := tok.(xml.StartElement); ok {
			return invalidXmlError
		}
		if t, ok := tok.(xml.EndElement); ok {
			if len(endElementName) == 0 {
				return nil
			}
			if t.Name.Local != endElementName {
				return invalidXmlError
			}
			return nil
		}
	}
}

func (dec *calldecoder) waitForEndElements(elementNames ...string) (err error) {
	for i := len(elementNames) - 1; i >= 0; i-- {
		err = dec.waitForEndElement(elementNames[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (dec *calldecoder) skipChildren() (err error) {
	level := 0
	var tok xml.Token
	for {
		if tok, err = dec.Token(); err != nil {
			return err
		}
		if _, ok := tok.(xml.StartElement); ok {
			level++
		}
		if _, ok := tok.(xml.EndElement); ok {
			if level == 0 {
				dec.nextToken = tok
				return nil
			}
			level--
		}
	}
}

func (dec *calldecoder) readNamedString(name string) (methodName string, err error) {
	_, err = dec.waitForStartElement(name)
	if err != nil {
		return "", err
	}

	data, err := dec.readCharData()
	if err != nil {
		return "", err
	}

	err = dec.waitForEndElement(name)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (dec *calldecoder) readMethodName() (methodName string, err error) {
	return dec.readNamedString("methodName")
}

// /////////////////////////////////////////////////////////////////////////////

type MethodCallParser interface {
	ParseMethodCall(methodName string, cb MethodCallParserCB) (err error)
}

// /////////////////////////////////////////////////////////////////////////////

type MethodCallParserCB interface {
	GetCallParam(val interface{}) (err error)
	IgnoreParams() (err error)
	PutResult(val interface{}) (err error)
}

type methodCallParserCB struct {
	dec      *calldecoder
	elements []string
}

func (cb *methodCallParserCB) GetCallParam(v interface{}) (err error) {
	_, err = cb.dec.requireStartElements(cb.elements...)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(v)
	err = cb.dec.decodeValue(val)
	if err != nil {
		return err
	}

	return cb.dec.waitForEndElements(cb.elements...)
}

func (cb *methodCallParserCB) IgnoreParams() (err error) {
	return cb.dec.skipChildren()
}

func (cb *methodCallParserCB) PutResult(v interface{}) (err error) {
	return nil
}

// /////////////////////////////////////////////////////////////////////////////

func newCallDecoder(data []byte) (dec *calldecoder) {
	dec = &calldecoder{decoder: decoder{xml.NewDecoder(bytes.NewBuffer(data))}}
	if CharsetReader != nil {
		dec.CharsetReader = CharsetReader
	}
	return dec
}

func (dec *calldecoder) parseMethodCall(mcp MethodCallParser) (err error) {
	_, err = dec.requireStartElements("methodCall")
	if err != nil {
		return err
	}
	methodName, err := dec.readMethodName()
	if err != nil {
		return err
	}
	_, err = dec.waitForStartElement("params")
	if err != nil {
		return err
	}

	if methodName != "system.multicall" {
		cb := &methodCallParserCB{
			dec:      dec,
			elements: []string{"param", "value"},
		}
		err = mcp.ParseMethodCall(methodName, cb)
		dec.result = append(dec.result, []byte("<value></value>")...)
		return err
	}

	_, err = dec.requireStartElements("param", "value", "array", "data")
	if err != nil {
		return err
	}

	dec.result = append(dec.result, []byte("<value><array><data>")...)

	for {
		tok, err := dec.waitForStartElements("value", "struct", "member")
		if err != nil {
			return err
		}
		if tok == nil {
			break
		}
		name, err := dec.readNamedString("name")
		if err != nil {
			return err
		}
		if name != "methodName" {
			return invalidXmlError
		}
		methodName, err := dec.readNamedString("value")
		if err != nil {
			return err
		}
		err = dec.waitForEndElement("member")
		if err != nil {
			return err
		}
		_, err = dec.requireStartElements("member")
		if err != nil {
			return err
		}
		name, err = dec.readNamedString("name")
		if err != nil {
			return err
		}
		if name != "params" {
			return invalidXmlError
		}
		_, err = dec.requireStartElements("value", "array", "data")
		if err != nil {
			return err
		}

		cb := &methodCallParserCB{
			dec:      dec,
			elements: []string{"value"},
		}
		err = mcp.ParseMethodCall(methodName, cb)
		if err != nil {
			return err
		}
		dec.result = append(dec.result, []byte("<value><array><data>")...)
		dec.result = append(dec.result, []byte("<value></value>")...)
		dec.result = append(dec.result, []byte("</data></array></value>")...)

		err = dec.waitForEndElements("value", "struct", "member", "value", "array", "data")
		if err != nil {
			return err
		}
	}

	dec.result = append(dec.result, []byte("</data></array></value>")...)

	err = dec.waitForEndElements("param", "value", "array", "data")
	if err != nil {
		return err
	}

	return nil
}

func HandleMethodCall(data []byte, mcp MethodCallParser) (response []byte, err error) {
	dec := newCallDecoder(data)
	err = dec.parseMethodCall(mcp)
	if err != nil {
		return nil, err
	}

	result := []byte("<?xml version=\"1.0\" encoding=\"iso-8859-1\"?><methodResponse><params><param>")
	result = append(result, dec.result...)
	result = append(result, []byte("</param></params></methodResponse>")...)

	return result, err
}
