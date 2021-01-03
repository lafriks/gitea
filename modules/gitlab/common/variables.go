// The MIT License (MIT)
// Copyright (c) 2015-2019 GitLab B.V.
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package common

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type JobVariable struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Public   bool   `json:"public"`
	Internal bool   `json:"-"`
	File     bool   `json:"file"`
	Masked   bool   `json:"masked"`
	Raw      bool   `json:"raw"`
}

type JobVariables []JobVariable

func (b JobVariable) String() string {
	return fmt.Sprintf("%s=%s", b.Key, b.Value)
}

func (b JobVariables) PublicOrInternal() (variables JobVariables) {
	for _, variable := range b {
		if variable.Public || variable.Internal {
			variables = append(variables, variable)
		}
	}
	return variables
}

func (b JobVariables) StringList() (variables []string) {
	for _, variable := range b {
		variables = append(variables, variable.String())
	}
	return variables
}

func (b JobVariables) Get(key string) string {
	switch key {
	case "$":
		return key
	case "*", "#", "@", "!", "?", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		return ""
	}
	for i := len(b) - 1; i >= 0; i-- {
		if b[i].Key == key {
			return b[i].Value
		}
	}
	return ""
}

// OverwriteKey overwrites an existing key with a new variable.
func (b JobVariables) OverwriteKey(key string, variable JobVariable) {
	for i, v := range b {
		if v.Key == key {
			b[i] = variable
			return
		}
	}
}

func (b JobVariables) ExpandValue(value string) string {
	return os.Expand(value, b.Get)
}

func (b JobVariables) Expand() JobVariables {
	var variables JobVariables
	for _, variable := range b {
		if !variable.Raw {
			variable.Value = b.ExpandValue(variable.Value)
		}

		variables = append(variables, variable)
	}
	return variables
}

func (b JobVariables) Masked() (masked []string) {
	for _, variable := range b {
		if variable.Masked {
			masked = append(masked, variable.Value)
		}
	}
	return
}

func ParseVariable(text string) (variable JobVariable, err error) {
	keyValue := strings.SplitN(text, "=", 2)
	if len(keyValue) != 2 {
		err = errors.New("missing =")
		return
	}
	variable = JobVariable{
		Key:   keyValue[0],
		Value: keyValue[1],
	}
	return
}
