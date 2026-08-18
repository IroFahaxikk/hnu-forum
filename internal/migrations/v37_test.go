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

	"github.com/apache/answer/internal/base/constant"
	"github.com/apache/answer/internal/entity"
	"github.com/stretchr/testify/require"
	"xorm.io/xorm"
)

func TestAddAnnouncementPopupReceipts(t *testing.T) {
	x, err := xorm.NewEngine("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() {
		_ = x.Close()
	}()

	require.NoError(t, addAnnouncementPopupReceipts(context.Background(), x))
	enabledAt := &entity.Config{Key: constant.AnnouncementPopupEnabledAtConfigKey}
	exists, err := x.Get(enabledAt)
	require.NoError(t, err)
	require.True(t, exists)
	require.NotEmpty(t, enabledAt.Value)
	_, err = x.Insert(&entity.AnnouncementReceipt{UserID: "1", QuestionID: "10"})
	require.NoError(t, err)
	_, err = x.Insert(&entity.AnnouncementReceipt{UserID: "1", QuestionID: "10"})
	require.Error(t, err)
}
