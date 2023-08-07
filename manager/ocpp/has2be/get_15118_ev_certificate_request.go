// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package has2be

import "encoding/json"
import "fmt"

type Get15118EVCertificateRequestJson struct {
	// Schema version currently used for the 15118 session between EV and Charging
	// Station. Needed for parsing of the EXI stream by the CSMS.
	//
	//
	A15118SchemaVersion *string `json:"15118SchemaVersion,omitempty" yaml:"15118SchemaVersion,omitempty" mapstructure:"15118SchemaVersion,omitempty"`

	// Raw CertificateInstallationReq request from EV, Base64 encoded.
	//
	ExiRequest string `json:"exiRequest" yaml:"exiRequest" mapstructure:"exiRequest"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Get15118EVCertificateRequestJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["exiRequest"]; !ok || v == nil {
		return fmt.Errorf("field exiRequest in Get15118EVCertificateRequestJson: required")
	}
	type Plain Get15118EVCertificateRequestJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Get15118EVCertificateRequestJson(plain)
	return nil
}
