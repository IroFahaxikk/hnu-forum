// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package content

import (
	"strings"
	"testing"

	"github.com/apache/answer/internal/schema"
	"github.com/stretchr/testify/assert"
)

func TestApplyGuestQuestionPreview(t *testing.T) {
	question := &schema.QuestionInfoResp{
		Content:        "## Preview heading\n\n**formatted preview** followed by private content",
		HTML:           "<p>complete body</p>",
		Description:    "plain-text fallback",
		MemberActions:  []*schema.PermissionMemberAction{{Name: "edit"}},
		ExtendsActions: []*schema.PermissionMemberAction{{Name: "invite"}},
	}

	applyGuestQuestionPreview(question)

	assert.Empty(t, question.Content)
	assert.Contains(t, question.HTML, `<h2 id="preview-heading">Preview heading</h2>`)
	assert.Contains(t, question.HTML, "<strong>formatted preview</strong>")
	assert.NotContains(t, question.HTML, "private content")
	assert.True(t, question.PreviewOnly)
	assert.Nil(t, question.MemberActions)
	assert.Nil(t, question.ExtendsActions)
}

func TestApplyGuestQuestionPreviewSanitizesFallback(t *testing.T) {
	question := &schema.QuestionInfoResp{
		Description: `preview <script>alert("xss")</script>`,
	}

	applyGuestQuestionPreview(question)

	assert.NotContains(t, question.HTML, "<script>")
	assert.True(t, question.PreviewOnly)
}

func TestBuildGuestQuestionPreviewCapsLongContent(t *testing.T) {
	content := strings.Repeat("内", guestQuestionPreviewMaxRunes+300)

	preview := buildGuestQuestionPreview(content)

	assert.Len(t, []rune(strings.TrimSuffix(preview, "\n\n…")), guestQuestionPreviewMaxRunes)
	assert.NotContains(t, preview, strings.Repeat("内", guestQuestionPreviewMaxRunes+1))
}

func TestApplyGuestQuestionPreviewWithNilQuestion(t *testing.T) {
	assert.NotPanics(t, func() {
		applyGuestQuestionPreview(nil)
	})
}
