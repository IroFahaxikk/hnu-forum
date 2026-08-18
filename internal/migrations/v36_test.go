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

package migrations

import (
	"context"
	"testing"
	"time"

	"github.com/apache/answer/internal/entity"
	"github.com/stretchr/testify/require"
	"xorm.io/xorm"
)

type questionBeforeFeatured struct {
	ID               string    `xorm:"not null pk BIGINT(20) id"`
	CreatedAt        time.Time `xorm:"not null default CURRENT_TIMESTAMP TIMESTAMP created_at"`
	UpdatedAt        time.Time `xorm:"updated_at TIMESTAMP"`
	UserID           string    `xorm:"not null default 0 BIGINT(20) INDEX user_id"`
	SectionID        int64     `xorm:"not null default 0 BIGINT(20) INDEX section_id"`
	InviteUserID     string    `xorm:"TEXT invite_user_id"`
	LastEditUserID   string    `xorm:"not null default 0 BIGINT(20) last_edit_user_id"`
	Title            string    `xorm:"not null default '' VARCHAR(150) title"`
	OriginalText     string    `xorm:"not null MEDIUMTEXT original_text"`
	ParsedText       string    `xorm:"not null MEDIUMTEXT parsed_text"`
	Pin              int       `xorm:"not null default 1 INT(11) pin"`
	Show             int       `xorm:"not null default 1 INT(11) show"`
	Status           int       `xorm:"not null default 1 INT(11) status"`
	ViewCount        int       `xorm:"not null default 0 INT(11) view_count"`
	UniqueViewCount  int       `xorm:"not null default 0 INT(11) unique_view_count"`
	VoteCount        int       `xorm:"not null default 0 INT(11) vote_count"`
	AnswerCount      int       `xorm:"not null default 0 INT(11) answer_count"`
	HotScore         int       `xorm:"not null default 0 INT(11) hot_score"`
	CollectionCount  int       `xorm:"not null default 0 INT(11) collection_count"`
	FollowCount      int       `xorm:"not null default 0 INT(11) follow_count"`
	AcceptedAnswerID string    `xorm:"not null default 0 BIGINT(20) accepted_answer_id"`
	LastAnswerID     string    `xorm:"not null default 0 BIGINT(20) last_answer_id"`
	PostUpdateTime   time.Time `xorm:"post_update_time TIMESTAMP"`
	RevisionID       string    `xorm:"not null default 0 BIGINT(20) revision_id"`
	LinkedCount      int       `xorm:"not null default 0 INT(11) linked_count"`
}

func (questionBeforeFeatured) TableName() string {
	return "question"
}

func TestAddFeaturedQuestions(t *testing.T) {
	x, err := xorm.NewEngine("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() {
		_ = x.Close()
	}()

	require.NoError(t, x.Sync(new(questionBeforeFeatured)))
	_, err = x.Insert(&questionBeforeFeatured{ID: "1"})
	require.NoError(t, err)

	require.NoError(t, addFeaturedQuestions(context.Background(), x))

	question := &entity.Question{}
	exists, err := x.ID("1").Get(question)
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, entity.QuestionUnfeatured, question.Featured)
	require.Zero(t, question.FeaturedAt)
}
