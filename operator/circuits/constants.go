package circuits

const (
	NOP = iota
	UPDATE
	INSERT
	DELETE
	zkeyFilePath           = "/circuit_final.zkey"
	wasmFilePath           = "/circuit.wasm"
	verficationKeyFilePath = "/verification_key.json"
)
