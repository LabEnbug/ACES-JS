"use client";

import Link from 'next/link'
import type { MenuProps } from 'antd';
import { Menu} from 'antd';
import React from 'react';


const items: MenuProps['items'] = [
  {
    label: (
      <Link
      href="/show">
        综合
      </Link>
    ),
    key: 'comprehensive',
  },
  {
    label: '娱乐',
    key: 'entertainment',
  },
  {
    label: '历史',
    key: 'history',
  },
  {
    label: '体育',
    key: 'physical',
  },
  {
    label: '美女',
    key: 'beauty',
  },
  {
    label: '读书',
    key: 'read',
  },
  {
    label: '搞笑',
    key: 'funny',
  },
  {
    label: '短剧场',
    key: 'shorttheater',
  },
  {
    label: '吃播',
    key: 'mukbang',
  },
  {
    label: (
      <a href="https://ant.design" target="_blank" rel="noopener noreferrer">
        生活
      </a>
    ),
    key: 'life',
  },
];

export default function Sider() {
  return (
      <div>
          <Menu
            mode="vertical"
            defaultSelectedKeys={['1']}
            defaultOpenKeys={['sub1']}
            style={{ height: '100%', borderRight: 0 }}
            items={items}
          />
      </div>
  )
}
