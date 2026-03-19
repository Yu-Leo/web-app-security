package models

type MLModel struct {
	ID        int64
	Name      string
	ModelData []byte
}

type MLModelToCreate struct {
	Name      string
	ModelData []byte
}

type MLModelToUpdate struct {
	ID        int64
	Name      string
	ModelData []byte
}
