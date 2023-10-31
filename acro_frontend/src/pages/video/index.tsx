import React from 'react';
import { Result, Button, Message} from '@arco-design/web-react';
import locale from './locale';
import useLocale from '@/utils/useLocale';
import styles from './style/index.module.less';
import VideoPlayer from '@/components/Video';
import { useState, useEffect } from 'react'
import axios from 'axios';
import { useRouter } from 'next/router';
import { useSelector, useDispatch } from 'react-redux';
import store, { GlobalState } from '@/store';
import { GlobalContext } from '@/context';
import GetAxios from '@/utils/getaxios';

function VideoP() {
  const t = useLocale(locale);
  const router = useRouter();
  const type = router.query['type'] || 'comprehensive';
  const video_uid = router.query['video_uid'];
  const pathname =  router.pathname;
  const [playlist, SetPlayList] = useState<any>([]);
  const [playindex, SetPlayIndex] = useState(0);
  const limit = 10;
  const page = 1;


  const JudgeStatus = (data: any)=> {
      if (data.status != 200) {
        Message.error(t['message.notfind'])
          return false
      }
      return true
  }

  const reflectplayindex = (index) => {
    SetPlayIndex(index)
    if (playlist.length > 0) {
      router.push({
        pathname: pathname,
        query: {
          'video_uid': playlist[index]['video_uid']
        },
      }, undefined, { shallow: true });
      window.localStorage.setItem('playvideo-id', playlist[index]['video_uid'])
    }
  }

  const recordhistory = () => {
    const pre = window.localStorage.getItem('playvideo-pre-id')
    const uid = window.localStorage.getItem('playvideo-id')
    if (uid && pre != uid) {
      const baxios  = GetAxios()
      const param = new FormData()
      param.append('video_uid', uid)
      baxios.post('/v1-api/v1/video/record', param)
      .then(response => {
        window.localStorage.setItem('playvideo-pre-id', uid)
      }).catch(error=>{
        console.error(error)
      })
    }
  }



  useEffect(() => {
    if (playlist.length == 0) {
      const baxios  = GetAxios()
      const param = new FormData()
      param.append('limit', limit)
      param.append('page', page)
      param.append('type', type)
      baxios.post('/v1-api/v1/video/list', param)
      .then(response => {
          const data = response.data
          if (JudgeStatus(data)) {
            if (video_uid && data.data.video_list.length >0 && video_uid != data.data.video_list[playindex]['video_uid']) {
              const param1 = new FormData()
              param1.append('video_uid', video_uid)
              baxios.post('v1-api/v1/video/info', param1).then(response1 => {
                if (JudgeStatus(response1.data)) {
                  data.data.video_list.unshift(response1.data.data.video)
                }
                SetPlayList(data.data.video_list)
              }).catch(error => {
                console.error(error)
              })
              // 
            } else SetPlayList(data.data.video_list)
          }
      })
      .catch(error => {
          console.error(error);
      });
    }
  }, [router.isReady, limit, page, type]);

  return (
    <div className={styles.container}>
      <div className={styles.wrapper}>
        <VideoPlayer  hlsPlayList={playlist} playindex={playindex} reflectplayindex={reflectplayindex} recordhistory={recordhistory} />
      </div>
    </div>
  );
}

export default VideoP;
