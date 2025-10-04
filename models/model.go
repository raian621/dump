package models

type ClientModel interface {
	ToStorageModel() StorageModel
}

type StorageModel interface {
	ToClientModel() ClientModel
}
