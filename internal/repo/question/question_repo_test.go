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

package question

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/apache/answer/internal/base/constant"
	"github.com/apache/answer/internal/base/data"
	"github.com/apache/answer/internal/entity"
	"github.com/apache/answer/internal/schema"
	"github.com/stretchr/testify/require"
	"xorm.io/xorm"
)

func TestGetQuestionPageFeaturedOnly(t *testing.T) {
	x, err := xorm.NewEngine("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() {
		_ = x.Close()
	}()
	require.NoError(t, x.Sync(new(entity.Question)))

	now := time.Now()
	questions := []*entity.Question{
		{ID: "1", CreatedAt: now.Add(-time.Hour), OriginalText: "one", ParsedText: "one", Featured: entity.QuestionFeatured, FeaturedAt: 10,
			Show: entity.QuestionShow, Status: entity.QuestionStatusAvailable},
		{ID: "2", CreatedAt: now, OriginalText: "two", ParsedText: "two", Featured: entity.QuestionUnfeatured,
			Show: entity.QuestionShow, Status: entity.QuestionStatusAvailable},
		{ID: "3", CreatedAt: now.Add(-2 * time.Hour), OriginalText: "three", ParsedText: "three", Featured: entity.QuestionFeatured, FeaturedAt: 20,
			Show: entity.QuestionShow, Status: entity.QuestionStatusAvailable},
	}
	_, err = x.Insert(questions)
	require.NoError(t, err)

	repo := &questionRepo{data: &data.Data{DB: x}}
	list, total, err := repo.GetQuestionPage(context.Background(), 1, 6, nil, nil, "",
		schema.QuestionOrderCondFeatured, 0, false, false)
	require.NoError(t, err)
	require.EqualValues(t, 2, total)
	require.Len(t, list, 2)
	require.Equal(t, "3", list[0].ID)
	require.Equal(t, "1", list[1].ID)
}

func TestClaimUnseenAnnouncements(t *testing.T) {
	x, err := xorm.NewEngine("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() {
		_ = x.Close()
	}()
	require.NoError(t, x.Sync(new(entity.Question), new(entity.AnnouncementReceipt), new(entity.Config)))

	now := time.Now()
	_, err = x.Insert(&entity.Config{Key: constant.AnnouncementPopupEnabledAtConfigKey,
		Value: fmt.Sprintf("%d", now.Add(-2*time.Hour).Unix())})
	require.NoError(t, err)
	questions := []*entity.Question{
		{ID: "10000000000000010", UserID: "admin", SectionID: entity.ForumSectionSiteAnnouncementsID,
			CreatedAt: now.Add(-time.Hour), OriginalText: "old", ParsedText: "old",
			Show: entity.QuestionShow, Status: entity.QuestionStatusAvailable},
		{ID: "10000000000000011", UserID: "admin", SectionID: entity.ForumSectionSiteAnnouncementsID,
			CreatedAt: now, OriginalText: "new", ParsedText: "new",
			Show: entity.QuestionShow, Status: entity.QuestionStatusAvailable},
		{ID: "10000000000000012", UserID: "user-1", SectionID: entity.ForumSectionSiteAnnouncementsID,
			CreatedAt: now, OriginalText: "own", ParsedText: "own",
			Show: entity.QuestionShow, Status: entity.QuestionStatusAvailable},
		{ID: "10000000000000013", UserID: "admin", SectionID: 101,
			CreatedAt: now, OriginalText: "normal", ParsedText: "normal",
			Show: entity.QuestionShow, Status: entity.QuestionStatusAvailable},
		{ID: "10000000000000014", UserID: "admin", SectionID: entity.ForumSectionSiteAnnouncementsID,
			CreatedAt: now.Add(-3 * time.Hour), OriginalText: "historical", ParsedText: "historical",
			Show: entity.QuestionShow, Status: entity.QuestionStatusAvailable},
	}
	_, err = x.Insert(questions)
	require.NoError(t, err)

	repo := &questionRepo{data: &data.Data{DB: x}}
	claimed, err := repo.ClaimUnseenAnnouncements(context.Background(), "user-1", 10)
	require.NoError(t, err)
	require.Len(t, claimed, 2)
	require.Equal(t, "10000000000000011", claimed[0].ID)
	require.Equal(t, "10000000000000010", claimed[1].ID)

	claimedAgain, err := repo.ClaimUnseenAnnouncements(context.Background(), "user-1", 10)
	require.NoError(t, err)
	require.Empty(t, claimedAgain)

	require.NoError(t, repo.MarkAnnouncementSeen(context.Background(), "user-2", "10000000000000011"))
	claimedAfterReading, err := repo.ClaimUnseenAnnouncements(context.Background(), "user-2", 10)
	require.NoError(t, err)
	require.Len(t, claimedAfterReading, 2)
	require.Equal(t, "10000000000000012", claimedAfterReading[0].ID)
	require.Equal(t, "10000000000000010", claimedAfterReading[1].ID)
}
