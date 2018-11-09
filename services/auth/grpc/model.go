package auth

type GRPCAuth interface {
	Add(string, string) (string, error)
	Exists(string) (bool, error)
}
