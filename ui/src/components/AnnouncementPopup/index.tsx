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

import { FC, useEffect, useRef, useState } from 'react';
import { ListGroup } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import Modal from '@/components/Modal';
import { pathFactory } from '@/router/pathFactory';
import { claimAnnouncementPopups } from '@/services';
import { loggedUserInfoStore } from '@/stores';
import type { AnnouncementPopupItem } from '@/common/interface';

const AnnouncementPopup: FC = () => {
  const { t } = useTranslation('translation', {
    keyPrefix: 'announcement_popup',
  });
  const { user } = loggedUserInfoStore();
  const requestedToken = useRef('');
  const [items, setItems] = useState<AnnouncementPopupItem[]>([]);

  useEffect(() => {
    if (!user.access_token) {
      requestedToken.current = '';
      setItems([]);
      return;
    }
    if (requestedToken.current === user.access_token) {
      return;
    }
    requestedToken.current = user.access_token;
    claimAnnouncementPopups()
      .then((announcements) => setItems(announcements || []))
      .catch(() => setItems([]));
  }, [user.access_token]);

  const close = () => setItems([]);

  return (
    <Modal
      visible={items.length > 0}
      title={t('title')}
      showConfirm={false}
      cancelText="close"
      scrollable
      onCancel={close}>
      <p className="text-secondary">{t('description')}</p>
      <ListGroup variant="flush">
        {items.map((item) => (
          <ListGroup.Item
            key={item.id}
            as={Link}
            action
            to={pathFactory.questionLanding(item.id, item.url_title)}
            onClick={close}>
            <div className="fw-semibold text-body">{item.title}</div>
            {item.description ? (
              <div className="small text-secondary mt-1 text-truncate-3">
                {item.description}
              </div>
            ) : null}
          </ListGroup.Item>
        ))}
      </ListGroup>
    </Modal>
  );
};

export default AnnouncementPopup;
