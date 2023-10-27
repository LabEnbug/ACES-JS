"use client";

import Link from 'next/link'
import MobileMenu from './mobile-menu'
import { useState, useEffect } from 'react'
import { Avatar, Button} from 'antd';
import axios from 'axios';
import type { MenuProps } from 'antd';
import { Input, Space } from 'antd';
const { Search } = Input;

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

export default function Header_() {
  const [showsigns, SetShowSigns] = useState<boolean>(true)
  const [nickname, SetNickName] = useState('')

  const loginout = (event: any) => {
    event.preventDefault();
    axios.get('/v1-api/v1/login')
    .then(response => {
      localStorage.removeItem("userInfo")
      SetShowSigns(true)
    })
    .catch(error => {
      console.error(error);
    });
  };

  useEffect(() => {
    const userInfo = window.localStorage.getItem('userInfo')
    if (userInfo) {
      const userinfo = JSON.parse(userInfo)
      if (userinfo.data.user.nickname && showsigns) {
        SetShowSigns(false)
        SetNickName(userinfo.data.user.nickname)
      }
    } 
  });
  return (
    <header className="absolute w-full z-30">
      <div className="max-w-6xl mx-auto px-4 sm:px-6">
        <div className="flex items-center justify-between h-20">
          {/* Site branding */}
          <div className="shrink-0 mr-4">
            {/* Logo */}
            <Link href="/" className="block" aria-label="Cruip">
              <svg className="w-8 h-8 fill-current text-purple-600" viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                <path d="M31.952 14.751a260.51 260.51 0 00-4.359-4.407C23.932 6.734 20.16 3.182 16.171 0c1.634.017 3.21.28 4.692.751 3.487 3.114 6.846 6.398 10.163 9.737.493 1.346.811 2.776.926 4.262zm-1.388 7.883c-2.496-2.597-5.051-5.12-7.737-7.471-3.706-3.246-10.693-9.81-15.736-7.418-4.552 2.158-4.717 10.543-4.96 16.238A15.926 15.926 0 010 16C0 9.799 3.528 4.421 8.686 1.766c1.82.593 3.593 1.675 5.038 2.587 6.569 4.14 12.29 9.71 17.792 15.57-.237.94-.557 1.846-.952 2.711zm-4.505 5.81a56.161 56.161 0 00-1.007-.823c-2.574-2.054-6.087-4.805-9.394-4.044-3.022.695-4.264 4.267-4.97 7.52a15.945 15.945 0 01-3.665-1.85c.366-3.242.89-6.675 2.405-9.364 2.315-4.107 6.287-3.072 9.613-1.132 3.36 1.96 6.417 4.572 9.313 7.417a16.097 16.097 0 01-2.295 2.275z" />
              </svg>
            </Link>
            </div>
            <div className="w-200 flex grow justify-end flex-wrap items-center ">
            <Search placeholder="搜索你感兴趣的内容"  enterButton />
            </div>
          {/* Desktop navigation */}
          <nav className="hidden md:flex md:grow">
            {showsigns ? <ul className="flex grow justify-end flex-wrap items-center">
              <li>
                <Link
                  href="/signin"
                  className="font-medium text-purple-600 hover:text-gray-200 px-4 py-3 flex items-center transition duration-150 ease-in-out"
                >
                  登录
                </Link>
              </li>
              <li>
                <Link href="/signup" className="btn-sm text-white bg-purple-600 hover:bg-purple-700 ml-3">
                  注册
                </Link>
              </li>
            </ul> : <ul className="flex grow justify-end flex-wrap items-center">
            <li>
              <Avatar src="https://xsgames.co/randomusers/avatar.php?g=pixel&key=1" />
            </li>
            <li>
              <Button type="link"> {nickname} </Button>
            </li>
            <li>
              <Button type="text" onClick={(e)=>{loginout(e)}}> 登出 </Button>
            </li>
            </ul>}
          </nav>
          <MobileMenu />
        </div>
      </div>
    </header>
  )
}
