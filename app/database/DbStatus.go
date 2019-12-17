package database

type DbStatus int32

const (
	DbSuccess DbStatus = iota
	DbFailed
	DbNotFound
	DbExisted
)
