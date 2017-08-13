package infrastructure

type Configuration struct {
	Server struct {
		Port int
	}
	Storage struct {
		Type string
		Path string
	}
}
