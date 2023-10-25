"use client";

import { useState, useEffect } from 'react'
import Link from 'next/link'
import axios from 'axios';
import  {useRouter } from "next/navigation"

export default function SignIn() {
  const [username, SetAccount] = useState('')
  const [password, SetPassward] = useState('')
  const router = useRouter()

  const handleSubmit = (event: any) => {
    event.preventDefault();
    const params = {
      username: 'user1',
      password: 'user1'
    }
    axios.get('/v1-api/v1/login', { params })
    .then(response => {
      const data = response.data
      localStorage.setItem('userInfo', JSON.stringify(data))
      router.refresh()
      console.log(data)
    })
    .catch(error => {
      console.error(error);
    });
  };

  useEffect(() => {
    const userInfo = window.localStorage.getItem('userInfo')
    if (userInfo) {
      router.push('/')
    } 
  });


  return (
    <section className="relative">
      <div className="max-w-6xl mx-auto px-4 sm:px-6">
        <div className="pt-32 pb-12 md:pt-40 md:pb-20">

          {/* Page header */}
          <div className="max-w-3xl text-black mx-auto text-center pb-12 md:pb-20">
            <h1 className="h1">欢迎回家，地球Online玩家！</h1>
          </div>

          {/* Form */}
          <div className="max-w-sm mx-auto">
            <form>
              <div className="flex flex-wrap -mx-3 mb-4">
                <div className="w-full px-3">
                  <label className="block text-gray-300 text-sm font-medium mb-1" htmlFor="email">账号</label>
                  <input id="email" type="email" className="form-input w-full text-gray-300" placeholder="邮箱" required value={username} onChange={(e)=>SetAccount(e.target.value)}  />
                </div>
              </div>
              <div className="flex flex-wrap -mx-3 mb-4">
                <div className="w-full px-3">
                  <label className="block text-gray-300 text-sm font-medium mb-1" htmlFor="password">密码</label>
                  <input id="password" type="password" className="form-input w-full text-gray-300" placeholder="**********" required value={password} onChange={(e)=>SetPassward(e.target.value)} />
                </div>
              </div>
              <div className="flex flex-wrap -mx-3 mb-4">
                <div className="w-full px-3">
                  <div className="flex justify-between">
                    <label className="flex items-center">
                      <input type="checkbox" className="form-checkbox" />
                      <span className="text-gray-400 ml-2">保持登录</span>
                    </label>
                    <Link href="/reset-password" className="text-purple-600 hover:text-gray-200 transition duration-150 ease-in-out">找回密码</Link>
                  </div>
                </div>
              </div>
              <div className="flex flex-wrap -mx-3 mt-6">
                <div className="w-full px-3">
                  <button className="btn text-white bg-purple-600 hover:bg-purple-700 w-full" onClick={(event) => handleSubmit(event)}>登录</button>
                </div>
              </div>
            </form>
            <div className="text-gray-400 text-center mt-6">
              想要加入我们？ <Link href="/signup" className="text-purple-600 hover:text-gray-200 transition duration-150 ease-in-out">注册</Link>
            </div>
          </div>

        </div>
      </div>
    </section>
  )
}
