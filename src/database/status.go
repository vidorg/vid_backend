package database

type DbStatus uint8

const (
	DbSuccess DbStatus = iota
	DbFailed
	DbNotFound
	DbExisted
)
