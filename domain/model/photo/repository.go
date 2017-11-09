package photo

type Repository interface {
	Save(photo Photo) (*Identifier, error)

	Read(id Identifier) (*Photo, error)

	Delete(id Identifier) error
}
