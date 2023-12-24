package xmlrpc

import (
	"testing"

	"golang.org/x/net/html/charset"
)

var listMethodsCall = `<?xml version="1.0" encoding="iso-8859-1"?>
<methodCall>
	<methodName>system.listMethods</methodName>
	<params>
		<param>
			<value><string>TESTCLIENT</string></value>
		</param>
	</params>
</methodCall>`

var multicall = `<?xml version="1.0" encoding="iso-8859-1"?>
<methodCall>
	<methodName>system.multicall</methodName>
	<params>
		<param>
			<value>
				<array>
					<data>
						<value>
							<struct>
								<member>
									<name>methodName</name>
									<value>event</value>
								</member>
								<member>
									<name>params</name>
									<value>
										<array>
											<data>
												<value>TESTCLIENT</value>
												<value>REQ0109197:1</value>
												<value>BOOT</value>
												<value><boolean>0</boolean></value>
											</data>
										</array>
									</value>
								</member>
							</struct>
						</value>
						<value>
							<struct>
								<member>
									<name>methodName</name>
									<value>event</value>
								</member>
								<member>
									<name>params</name>
									<value>
										<array>
											<data>
												<value>TESTCLIENT</value>
												<value>REQ0109197:1</value>
												<value>GAS_ENERGY_COUNTER</value>
												<value><double>13.870000</double></value>
											</data>
										</array>
									</value>
								</member>
							</struct>
						</value>
						<value>
							<struct>
								<member>
									<name>methodName</name>
									<value>event</value>
								</member>
								<member>
									<name>params</name>
									<value>
										<array>
											<data>
												<value>TESTCLIENT</value>
												<value>REQ0109197:1</value>
												<value>GAS_POWER</value>
												<value><double>2.250000</double></value>
											</data>
										</array>
									</value>
								</member>
							</struct>
						</value>
					</data>
				</array>
			</value>
		</param>
	</params>
</methodCall>`

type checkMethodName struct {
	methodNames []string
}

func (cmn *checkMethodName) ParseMethodCall(methodName string, cb MethodCallParserCB) (err error) {
	cmn.methodNames = append(cmn.methodNames, methodName)
	switch methodName {
	case "system.listMethods":
		var clientName interface{}
		err = cb.GetCallParam(&clientName)
		if err != nil {
			return err
		}
	case "event":
		var interfaceID string
		var address string
		var valueKey string
		var value interface{}
		err = cb.GetCallParam(&interfaceID)
		if err != nil {
			return err
		}
		err = cb.GetCallParam(&address)
		if err != nil {
			return err
		}
		err = cb.GetCallParam(&valueKey)
		if err != nil {
			return err
		}
		err = cb.GetCallParam(&value)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestCallDecoderListMethods(t *testing.T) {
	CharsetReader = charset.NewReaderLabel

	cmn := &checkMethodName{}

	if _, err := HandleMethodCall([]byte(listMethodsCall), cmn); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
}

func TestCallDecoderMultiCall(t *testing.T) {
	CharsetReader = charset.NewReaderLabel

	cmn := &checkMethodName{}

	if _, err := HandleMethodCall([]byte(multicall), cmn); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
}
