package blocks

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Block interface {
	Render() string
}

type FormData struct {
	Index  int
	ID     string
	Fields []struct {
		Label string
		Kind  string
		Value string
	}
}

type Blocks struct {
	Blocks map[string]Block
}

func New() Blocks {
	return Blocks{
		Blocks: map[string]Block{},
	}
}

func (b *Blocks) RegisterBlock(id string, block Block) error {
	if _, ok := b.Blocks[id]; ok {
		return errors.New("block already exists")
	}

	b.Blocks[id] = block

	return nil
}

func (b *Blocks) GetRegisteredBlocks() []string {
	var blocks []string
	for id := range b.Blocks {
		blocks = append(blocks, id)
	}
	return blocks
}

func (b *Blocks) GetRegisteredBlock(id string) (Block, error) {
	block, ok := b.Blocks[id]
	if !ok {
		return nil, errors.New("block not found")
	}

	return block, nil
}

func (b *Blocks) GetFormDataByID(id string) (FormData, error) {
	form := FormData{
		Index: 0,
		ID:    id,
		Fields: make([]struct {
			Label string
			Kind  string
			Value string
		}, 0),
	}

	block, ok := b.Blocks[id]
	if !ok {
		return form, errors.New("block not found")
	}

	val := reflect.ValueOf(block)

	if val.Kind() != reflect.Struct {
		return form, errors.New("block is not a struct")
	}

	for i := range val.NumField() {
		label := val.Type().Field(i).Name
		kind := ""
		value := ""

		switch val.Field(i).Kind() {
		case reflect.String:
			kind = "string"
		case reflect.Int:
			kind = "number"
		case reflect.Bool:
			kind = "bool"
		default:
			fmt.Printf("Ignoring unsupported type %d on %s\n", val.Field(i).Kind(), label)
			continue
		}

		form.Fields = append(form.Fields, struct {
			Label string
			Kind  string
			Value string
		}{
			Label: label,
			Kind:  kind,
			Value: value,
		})
	}

	return form, nil
}

func (b *Blocks) GetFormDataByType(block Block) (FormData, error) {
	form := FormData{
		Index: 0,
		ID:    "",
		Fields: make([]struct {
			Label string
			Kind  string
			Value string
		}, 0),
	}

	if block == nil {
		return form, fmt.Errorf("block is nil")
	}

	val := reflect.ValueOf(block)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return form, fmt.Errorf("block pointer is nil")
		}

		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return form, fmt.Errorf("block is not a struct or pointer to struct, got %s", val.Kind())
	}

	for i := range val.NumField() {
		label := val.Type().Field(i).Name
		kind := ""
		value := val.Field(i).String()

		switch val.Field(i).Kind() {
		case reflect.String:
			kind = "string"
		case reflect.Int:
			kind = "number"
		case reflect.Bool:
			kind = "bool"
		default:
			fmt.Printf("Ignoring unsupported type %d on %s\n", val.Field(i).Kind(), label)
			continue
		}

		form.Fields = append(form.Fields, struct {
			Label string
			Kind  string
			Value string
		}{
			Label: label,
			Kind:  kind,
			Value: value,
		})
	}

	return form, nil
}

func (b *Blocks) ParseForm(id string, fields map[string]string) (Block, error) {
	formBlock, err := b.GetRegisteredBlock(id)
	if err != nil {
		return nil, err
	}

	blockType := reflect.TypeOf(formBlock)
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
			continue
		}
	}

	result := reflect.New(blockType).Elem()
	result.Set(blockValue)
	return result.Addr().Interface().(Block), nil
}
