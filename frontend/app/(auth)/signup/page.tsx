"use client";

import { useState, useEffect } from 'react'
import Link from 'next/link'
import axios from 'axios';
import  {useRouter } from "next/navigation"
import { EyeInvisibleOutlined, EyeTwoTone } from '@ant-design/icons';
import { Input, message } from 'antd';

export default function SignUp() {
  const [username, SetAccount] = useState('')
  const [password, SetPassward] = useState('')
  const [confirmpassword, SetConfirmPassward] = useState('')
  const [cpstatus, SetCpStatus] = useState<any>("warning")
  const [nickname, SetNickname] = useState('')
  const router = useRouter()

  useEffect(() => {
    const userInfo = window.localStorage.getItem('userInfo')
    if (userInfo) {
      router.push('/')
    } 
  });

  const onConfirmPasswardChange = (value: any) => {
    if (password == value) {
      SetCpStatus(null)
    }
    SetConfirmPassward(value)
  }; 

  const handleSubmit = (event: any) => {
    event.preventDefault();
    if (confirmpassword != password) {
      message.error("两次密码输入不一致！")
      return
    }
    const params = {
      username,
      password,
      nickname
    }
    axios.get('/v1-api/v1/signup', { params })
    .then(response => {
      const data = response.data
      if (data.status != 200) {
        message.error("邮箱已存在！")
      } else {
        message.info("注册成功，准备跳转首页", 3, onclose=()=>{
          router.refresh()
        })
        localStorage.setItem('userInfo', JSON.stringify(data))
      }
    })
    .catch(error => {
      console.error(error);
    });
  };

  return (
    <section className="relative">
      <div className="max-w-6xl mx-auto px-4 sm:px-6">
        <div className="pt-32 pb-12 md:pt-40 md:pb-20">

          {/* Page header */}
          <div className="max-w-3xl mx-auto text-black text-center pb-12 md:pb-20">
            <h1 className="h1">开始你的绚烂旅程！</h1>
          </div>
          {/* Form */}
          <div className="max-w-sm mx-auto">
            <form>
              <div className="flex flex-wrap -mx-3 mb-4">
                <div className="w-full px-3">
                  <label className="block text-gray-300 text-sm font-medium mb-1" htmlFor="full-name">昵称<span className="text-red-600">*</span></label>
                  <Input size="large" placeholder="一个属于你的名字" required value={nickname} onChange={(e)=>{SetNickname(e.target.value)}}/>
                </div>
              </div>
              <div className="flex flex-wrap -mx-3 mb-4">
                <div className="w-full px-3">
                  <label className="block text-gray-300 text-sm font-medium mb-1" htmlFor="email">邮箱</label>
                  <Input placeholder="一个属于你的邮箱"  size="large" type="email" required value={username} onChange={(e)=>SetAccount(e.target.value)} />
                </div>
              </div>
              <div className="flex flex-wrap -mx-3 mb-4">
                <div className="w-full px-3">
                  <label className="block text-gray-300 text-sm font-medium mb-1" htmlFor="password">密码 <span className="text-red-600">*</span></label>
                  <Input.Password size="large" placeholder="一个属于你的密码" required onChange={(e)=>{SetPassward(e.target.value)}} iconRender={(visible) => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)} />
                </div>
              </div>
              <div className="flex flex-wrap -mx-3 mb-4">
                <div className="w-full px-3">
                  <label className="block text-gray-300 text-sm font-medium mb-1" htmlFor="password">确认密码 <span className="text-red-600">*</span></label>
                  <Input.Password size="large" status={cpstatus}  placeholder="重新输入一遍你的密码" required onChange={(e)=>{onConfirmPasswardChange(e.target.value)}} iconRender={(visible) => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)} />
                </div>
              </div>
              <div className="flex flex-wrap -mx-3 mt-6">
                <div className="w-full px-3">
                  <button className="btn text-white bg-purple-600 hover:bg-purple-700 w-full" onClick={(e)=> {handleSubmit(e)}}>Sign up</button>
                </div>
              </div>
            </form>
            <div className="text-gray-400 text-center mt-6">
              Already using Open PRO? <Link href="/signin" className="text-purple-600 hover:text-gray-200 transition duration-150 ease-in-out">Sign in</Link>
            </div>
          </div>

        </div>
      </div>
    </section>
  )
}
