package blocks

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/google/uuid"
)

type Block interface {
	Render() string
}

type Handle struct {
	ID   string
	Type string
	Data Block
}

type FormField struct {
	Label string
	Value string
	Kind  string
}

type FormData struct {
	ID     string
	Name   string
	Fields []FormField
}

type Blocks struct {
	Store   map[string]Block
	Handles map[string]Handle
}

func New() Blocks {
	return Blocks{
		Store:   map[string]Block{},
		Handles: map[string]Handle{},
	}
}

func (b Blocks) RegisterBlock(name string, block Block) error {
	if _, ok := b.Store[name]; ok {
		return errors.New("block already exists")
	}

	b.Store[name] = block

	return nil
}

func (b Blocks) GetRegisteredBlocksIDs() []string {
	var blocks []string

	for id := range b.Store {
		blocks = append(blocks, id)
	}

	return blocks
}

func (b Blocks) NewBlock(name string) (Handle, error) {
	id := uuid.New().String()

	block, ok := b.Store[name]
	if !ok {
		return Handle{}, errors.New("block does not exist")
	}

	b.Handles[id] = Handle{
		ID:   id,
		Type: name,
		Data: block,
	}

	return b.Handles[id], nil
}

func (b Blocks) UpdateBlock(handle Handle) error {
	if _, ok := b.Handles[handle.ID]; !ok {
		return errors.New("cannot find block")
	}

	b.Handles[handle.ID] = handle

	return nil
}

func (b Blocks) GetBlock(id string) (Handle, error) {
	handle, ok := b.Handles[id]
	if !ok {
		return Handle{}, errors.New("block not found")
	}

	return handle, nil
}

func (b Blocks) DeleteBlock(handle Handle) error {
	delete(b.Handles, handle.ID)

	return nil
}

func (b Blocks) GetFormData(handle Handle) (FormData, error) {
	if _, ok := b.Handles[handle.ID]; !ok {
		return FormData{}, errors.New("block not found")
	}

	form := FormData{
		ID:     handle.ID,
		Name:   handle.Type,
		Fields: []FormField{},
	}

	val := reflect.ValueOf(handle.Data)

	if val.Kind() != reflect.Struct {
		return form, errors.New("block is not a struct")
	}

	for i := range val.NumField() {
		field := FormField{
			Label: val.Type().Field(i).Name,
			Kind:  "",
			Value: val.Field(i).String(),
		}

		switch val.Field(i).Kind() {
		case reflect.String:
			field.Kind = "string"
		case reflect.Int:
			field.Kind = "number"
		case reflect.Bool:
			field.Kind = "bool"
		default:
			fmt.Printf("Ignoring unsupported type %d on %s\n", val.Field(i).Kind(), field.Label)
			continue
		}

		form.Fields = append(form.Fields, field)
	}

	return form, nil
}

func (b Blocks) ParseFormIntoBlock(fields map[string]string, handle Handle) (Handle, error) {
	if _, ok := b.Handles[handle.ID]; !ok {
		return Handle{}, errors.New("block not found")
	}

	blockType := reflect.TypeOf(handle.Data)
	blockValue := reflect.New(blockType).Elem()
	for i := range blockValue.NumField() {
		label := blockValue.Type().Field(i).Name

		switch blockValue.Field(i).Kind() {
		case reflect.String:
			value, ok := fields[label]
			if !ok {
				blockValue.Field(i).SetString("")
			}

			blockValue.Field(i).SetString(value)
		case reflect.Int:
			value, ok := fields[label]
			if !ok {
				blockValue.Field(i).SetInt(0)
			}

			atoi, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				blockValue.Field(i).SetInt(0)
			}

			blockValue.Field(i).SetInt(atoi)
		case reflect.Bool:
			value, ok := fields[label]
			if !ok {
				blockValue.Field(i).SetBool(false)
			}

			if value == "on" {
				blockValue.Field(i).SetBool(true)
			} else {
				blockValue.Field(i).SetBool(false)
			}
		default:
			fmt.Printf("Ignoring unsupported type %d on %s\n", blockValue.Field(i).Kind(), label)
		}
	}

	result := reflect.New(blockType).Elem()
	result.Set(blockValue)

	handle.Data = result.Addr().Elem().Interface().(Block)

	b.Handles[handle.ID] = handle

	return handle, nil
}

func (b Blocks) Render(handle Handle) string {
	html := handle.Data.Render()

	return html
}
