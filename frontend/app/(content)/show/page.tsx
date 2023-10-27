'use client'

import axios from 'axios';
import { useState, useEffect } from 'react'
import { message } from 'antd';
import {VideoPlayer} from '@/components/short-video'
// https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8

export default function VideoPage() {
    const [playlist, SetPlayList] = useState<any>([])
    const [playindex, SetPlayIndex] = useState(0)
    const limit = 10
    const page = 1

    const JudgeStatus = (data: any)=> {
        if (data.status != 200) {
            message.error("找不到视频")
            return false
        }
        return true
    }

    useEffect(() => {
        const userInfo = window.localStorage.getItem('userInfo')
        console.log(userInfo)
        if (userInfo) {
          const userinfo = JSON.parse(userInfo)
          
        } else if (playlist.length == 0) {
            const params = {
                limit,
                page
            }
            axios.get('/v1-api/v1/video/list', {params})
            .then(response => {
                const data = response.data
                console.log(data)
                if (JudgeStatus(data)) {
                    SetPlayList(data.data.video_list)
                }
            })
            .catch(error => {
                console.error(error);
            });
        }
    });
    return (
        <div className="mt-[10px] bg-cover bg-black-500 w-full h-full">
            <div className='h-full w-full'>
                <VideoPlayer  hlsPlayList={playlist} playindex={playindex} />
            </div>
        </div>
        
    )
}
  