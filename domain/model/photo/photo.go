package photo

type Photo struct {
	id    Identifier
	image []byte
}

func New(data []byte) *Photo {
	return &Photo{Identifier{}, data}
}

func Of(id Identifier, data []byte) *Photo {
	return &Photo{id, data}
}

func (photo *Photo) Image() []byte {
	return photo.image
}

func (photo *Photo) Id() Identifier {
	return photo.id
}

func (photo *Photo) IsNew() bool {
	return len(photo.id.value) == 0
}
