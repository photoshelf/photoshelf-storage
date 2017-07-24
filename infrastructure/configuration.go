package infrastructure

type Configuration struct {
	Server struct {
		Port int
	}
	Storage struct {
		Directory string
	}
}
