import React from 'react';
import { Result, Button, Message} from '@arco-design/web-react';
import locale from './locale';
import useLocale from '@/utils/useLocale';
import styles from './style/index.module.less';
import VideoPlayer from '@/components/Video';
import { useState, useEffect } from 'react'
import axios from 'axios';
import { useSelector, useDispatch } from 'react-redux';
import store, { GlobalState } from '@/store';
import { GlobalContext } from '@/context';

function Comprehensive() {
  const { baxios} = useSelector((state: GlobalState) => state);
  const t = useLocale(locale);
  const [playlist, SetPlayList] = useState<any>([])
  const [playindex, SetPlayIndex] = useState(0)
  const limit = 10
  const page = 1

  const JudgeStatus = (data: any)=> {
      if (data.status != 200) {
        Message.error(t['message.notfind'])
          return false
      }
      return true
  }

  useEffect(() => {
      const userInfo = window.localStorage.getItem('userInfo')
      
      if (playlist.length == 0) {
          const params = {
              limit,
              page
          }
          const instance = axios.create({
            headers: {
              'Referrer Policy': 'origin',
            }
          });
          if (baxios) {
            baxios.get('/v1-api/v1/video/list', {params})
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
      }
  });
  return (
    <div className={styles.container}>
      <div className={styles.wrapper}>
        {/* <Result
          className={styles.result}
          status="403"
          subTitle={t['exception.result.403.description']}
          extra={
            <Button key="back" type="primary">
              {t['exception.result.403.back']}
            </Button>
          }
        /> */}
        <VideoPlayer  hlsPlayList={playlist} playindex={playindex} />
      </div>
    </div>
  );
}

export default Comprehensive;
