/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package schema

import (
	"testing"

	basevalidator "github.com/apache/answer/internal/base/validator"
	"github.com/segmentfault/pacman/i18n"
	"github.com/stretchr/testify/require"
)

func TestPostContentOnlyRequiresNonBlankText(t *testing.T) {
	tests := []struct {
		name        string
		request     any
		expectError bool
	}{
		{
			name:    "one character question",
			request: &QuestionAdd{Title: "题", Content: "文", SectionID: 1},
		},
		{
			name:        "blank question title",
			request:     &QuestionAdd{Title: " \n", Content: "文", SectionID: 1},
			expectError: true,
		},
		{
			name:        "blank question content",
			request:     &QuestionAdd{Title: "题", Content: "\t\n", SectionID: 1},
			expectError: true,
		},
		{
			name:    "one character answer",
			request: &AnswerAddReq{QuestionID: "1", Content: "答"},
		},
		{
			name:        "blank answer",
			request:     &AnswerAddReq{QuestionID: "1", Content: "  "},
			expectError: true,
		},
		{
			name:    "one character reply",
			request: &AddCommentReq{ObjectID: "1", OriginalText: "回"},
		},
		{
			name:        "blank reply",
			request:     &AddCommentReq{ObjectID: "1", OriginalText: " \t"},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := basevalidator.GetValidatorByLang(i18n.DefaultLanguage).Check(test.request)
			if test.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
