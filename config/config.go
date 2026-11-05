package config

type Config struct {
	CharacterName string
	AuthToken     string
}

func NewConfig() Config {
	return Config{
		CharacterName: "Buba",
		AuthToken:     "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6InRpbXVyMzMzOTlAdmsuY29tIiwicGFzc3dvcmRfY2hhbmdlZCI6bnVsbH0.qak98F18TUDBY8CLiK1ntzcD8SRGiGDhsgS916A1N1U",
	}
}
