package config


type CLIArgs struct {
	ModelId   string
	Question  string
	Provider  string
	Configure bool
}


type Config struct {
	ModelId  string `json:"ModelId"`
	Provider string `json:"Provider"`
}
