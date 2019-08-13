package models

// implement this interface if model needs to receive as json
type BaseModel interface {
	CheckValid() bool
}
