package option

// type ShadowTLSInboundOptions struct {
// 	ListenOptions
// 	Version                int                                                  `json:"version,omitempty"`
// 	Password               string                                               `json:"password,omitempty"`
// 	Users                  []ShadowTLSUser                                      `json:"users,omitempty"`
// 	Handshake              ShadowTLSHandshakeOptions                            `json:"handshake,omitempty"`
// 	HandshakeForServerName *badjson.TypedMap[string, ShadowTLSHandshakeOptions] `json:"handshake_for_server_name,omitempty"`
// 	StrictMode             bool                                                 `json:"strict_mode,omitempty"`
// 	WildcardSNI            WildcardSNI                                          `json:"wildcard_sni,omitempty"`
// }

type PPTPOutboundOptions struct {
	DialerOptions
	ServerOptions
	Username string `json:"username"`
	Password string `json:"password"`
}
