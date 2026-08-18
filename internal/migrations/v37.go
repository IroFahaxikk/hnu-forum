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
	"fmt"
	"strconv"
	"time"

	"github.com/apache/answer/internal/base/constant"
	"github.com/apache/answer/internal/entity"
	"xorm.io/xorm"
)

func addAnnouncementPopupReceipts(ctx context.Context, x *xorm.Engine) error {
	if err := x.Context(ctx).Sync(new(entity.AnnouncementReceipt), new(entity.Config)); err != nil {
		return fmt.Errorf("sync announcement receipt table failed: %w", err)
	}
	enabledAt := &entity.Config{Key: constant.AnnouncementPopupEnabledAtConfigKey}
	exists, err := x.Context(ctx).Get(enabledAt)
	if err != nil {
		return fmt.Errorf("get announcement popup activation time failed: %w", err)
	}
	if !exists {
		enabledAt.Value = strconv.FormatInt(time.Now().Unix(), 10)
		if _, err = x.Context(ctx).Insert(enabledAt); err != nil {
			return fmt.Errorf("save announcement popup activation time failed: %w", err)
		}
	}
	return nil
}
