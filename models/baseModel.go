package models

// implement this interface
// if model needs to receive as json (check invalid)
// or query in database (check not exist)
type BaseModel interface {
	CheckValid() bool
}
