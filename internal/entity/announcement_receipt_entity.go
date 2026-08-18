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

package entity

import "time"

// AnnouncementReceipt records that a user has either opened an announcement
// or had it delivered in the announcement popup.
type AnnouncementReceipt struct {
	UserID     string    `xorm:"not null pk BIGINT(20) user_id"`
	QuestionID string    `xorm:"not null pk BIGINT(20) question_id"`
	SeenAt     time.Time `xorm:"not null default CURRENT_TIMESTAMP TIMESTAMP seen_at"`
}

func (AnnouncementReceipt) TableName() string {
	return "announcement_receipt"
}
