package ocpi

var (
	// success
	StatusSuccess = int32(1000)
	// client errors
	StatusGenericClientError   = int32(2000)
	StatusInvalidParameters    = int32(2001)
	StatusNotEnoughInformation = int32(2002)
	StatusUnknownLocation      = int32(2003)
	StatusUnknownToken         = int32(2004)
	// server errors
	StatusGenericServerFailure = int32(3000)
	StatusUnableToUseClientApi = int32(3001)
	StatusUnsupportedVersion   = int32(3002)
	StatusNoMatchingEndpoints  = int32(3003)
	// hub errors
	StatusUnknownReceiver = int32(4001)
	StatusTimeout         = int32(4002)
	StatusConnectionError = int32(4003)
)

var (
	// success
	StatusSuccessMessage = "Success"
	// client errors
	StatusGenericClientErrorMessage   = "Generic client error"
	StatusInvalidParametersMessage    = "Invalid parameters"
	StatusNotEnoughInformationMessage = "Not enough information"
	StatusUnknownLocationMessage      = "Unknown location"
	StatusUnknownTokenMessage         = "Unknown token"
	// server errors
	StatusGenericServerFailureMessage = "Generic server failure"
	StatusUnableToUseClientApiMessage = "Unable to use client API"
	StatusUnsupportedVersionMessage   = "Unsupported version"
	StatusNoMatchingEndpointsMessage  = "No matching endpoints"
	// hub errors
	StatusUnknownReceiverMessage = "Unknown receiver"
	StatusTimeoutMessage         = "Timeout"
	StatusConnectionErrorMessage = "Connection error"
	// other
	StatusOtherMessage = "Other"
)
