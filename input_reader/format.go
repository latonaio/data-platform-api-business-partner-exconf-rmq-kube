package input_reader

import (
	"encoding/json"

	"golang.org/x/xerrors"
)

func ConvertToInput(data map[string]interface{}) (*Input, error) {
	inputString, err := json.Marshal(data)
	if err != nil {
		return nil, xerrors.Errorf("unknown error: %w", err)
	}
	input := &Input{}
	err = json.Unmarshal(inputString, input)
	if err != nil {
		return nil, xerrors.Errorf("cannot convert input file to json: %w", err)
	}
	return input, nil
}
