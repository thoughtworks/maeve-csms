package ocpp201

type TransactionEventResponseJson struct {
	// Priority from a business point of view. Default priority is 0, The range is
	// from -9 to 9. Higher values indicate a higher priority. The chargingPriority in
	// &lt;&lt;transactioneventresponse,TransactionEventResponse&gt;&gt; is
	// temporarily, so it may not be set in the
	// &lt;&lt;cmn_idtokeninfotype,IdTokenInfoType&gt;&gt; afterwards. Also the
	// chargingPriority in
	// &lt;&lt;transactioneventresponse,TransactionEventResponse&gt;&gt; overrules the
	// one in &lt;&lt;cmn_idtokeninfotype,IdTokenInfoType&gt;&gt;.
	//
	ChargingPriority *int `json:"chargingPriority,omitempty" yaml:"chargingPriority,omitempty" mapstructure:"chargingPriority,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// IdTokenInfo corresponds to the JSON schema field "idTokenInfo".
	IdTokenInfo *IdTokenInfoType `json:"idTokenInfo,omitempty" yaml:"idTokenInfo,omitempty" mapstructure:"idTokenInfo,omitempty"`

	// SHALL only be sent when charging has ended. Final total cost of this
	// transaction, including taxes. In the currency configured with the Configuration
	// Variable: &lt;&lt;configkey-currency,`Currency`&gt;&gt;. When omitted, the
	// transaction was NOT free. To indicate a free transaction, the CSMS SHALL send
	// 0.00.
	//
	//
	TotalCost *float64 `json:"totalCost,omitempty" yaml:"totalCost,omitempty" mapstructure:"totalCost,omitempty"`

	// UpdatedPersonalMessage corresponds to the JSON schema field
	// "updatedPersonalMessage".
	UpdatedPersonalMessage *MessageContentType `json:"updatedPersonalMessage,omitempty" yaml:"updatedPersonalMessage,omitempty" mapstructure:"updatedPersonalMessage,omitempty"`
}

func (*TransactionEventResponseJson) IsResponse() {}
