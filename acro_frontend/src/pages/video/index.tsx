import React from 'react';
import { Message } from '@arco-design/web-react';
import locale from './locale';
import useLocale from '@/utils/useLocale';
import styles from './style/index.module.less';
import VideoPlayer from '@/components/Video';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import GetVideType from '@/utils/getvideotype';
import baxios from "@/utils/getaxios";

function VideoP() {
  const default_type = 'comprehensive';
  const t = useLocale(locale);
  const router = useRouter();
  const type = router.query['type'] || default_type;
  const video_uid = router.query['video_uid'];
  const pathname = router.pathname;
  const [playlist, SetPlayList] = useState<any>([]);
  const [playIndex, SetPlayIndex] = useState(0);
  const limit = 10;

  // const page = 1;

  const JudgeStatus = (data: any) => {
    if (data.status != 200) {
      Message.error(t['message.notfind']);
      return false;
    }
    return true;
  };

  const reflectPlayIndex = (index) => {
    SetPlayIndex(index);
    const pre_type = window.localStorage.getItem('pretype') || 'comprehensive';
    if (playlist.length > 0) {
      router.push(
        {
          pathname: pathname,
          query: {
            type: pre_type,
            video_uid: playlist[index]['video_uid'],
          },
        },
        undefined,
        { shallow: true }
      );
      window.localStorage.setItem('playvideo-id', playlist[index]['video_uid']);
    }
  };

  const recordWatched = () => {
    const pre = window.localStorage.getItem('playvideo-pre-id');
    const uid = window.localStorage.getItem('playvideo-id');
    if (uid && pre != uid) {
      const param = new FormData();
      param.append('video_uid', uid);
      baxios
        .post('/v1-api/v1/video/watch', param)
        .then((response) => {
          window.localStorage.setItem('playvideo-pre-id', uid);
        })
        .catch((error) => {
          console.error(error);
        });
    }
  };

  useEffect(() => {
    const pre_type = window.localStorage.getItem('pretype');

    if (type != default_type && GetVideType(type) >= 999) {
      router.push('/video');
    }
    if (playlist.length == 0 || pre_type != type) {
      const param = new FormData();
      param.append('limit', limit);
      // param.append('page', page)
      if (type != default_type) param.append('type', GetVideType(type));
      baxios
        .post('/v1-api/v1/video/list', param)
        .then((response) => {
          const data = response.data;
          window.localStorage.setItem('pretype', type);
          if (JudgeStatus(data)) {
            if (
              video_uid &&
              data.data.video_list.length > 0 &&
              video_uid != data.data.video_list[playIndex]['video_uid']
            ) {
              const param1 = new FormData();
              param1.append('video_uid', video_uid);
              baxios
                .post('v1-api/v1/video/info', param1)
                .then((response1) => {
                  if (JudgeStatus(response1.data)) {
                    data.data.video_list.unshift(response1.data.data.video);
                  }
                  SetPlayList(data.data.video_list);
                })
                .catch((error) => {
                  console.error(error);
                });
              //
            } else SetPlayList(data.data.video_list);
          }
        })
        .catch((error) => {
          console.error(error);
        });
    } else if (playIndex >= playlist.length - 3) {
      const param = new FormData();
      param.append('limit', limit);
      if (type != default_type) param.append('type', GetVideType(type));
      param.append('start', playlist.length);
      baxios
        .post('/v1-api/v1/video/list', param)
        .then((res) => {
          const data = res.data;
          if (JudgeStatus(data)) {
            SetPlayList(playlist.concat(data.data.video_list));
          }
        })
        .catch((error) => {
          console.error(error);
        });
    }
  }, [router.isReady, type, video_uid, playIndex]);

  useEffect(() => {
    const lastHistoryState = window.history.state;

    const handleRouteChange = (url) => {
      // 检查新的state是否小于旧的state，如果是，那么很可能是后退操作
      const currentHistoryState = window.history.state;
      if (currentHistoryState.idx < lastHistoryState.idx) {
        window.location.reload();
        // 这里你可以添加你的业务逻辑
      }
    };

    router.events.on('routeChangeStart', handleRouteChange);

    return () => {
      router.events.off('routeChangeStart', handleRouteChange);
    };
  }, [router]);

  return (
    <div className={styles.container}>
      <div className={styles.wrapper}>
        <VideoPlayer
          hlsPlayList={playlist}
          playIndex={playIndex}
          reflectPlayIndex={reflectPlayIndex}
          recordWatched={recordWatched}
          options={undefined}
        />
      </div>
    </div>
  );
}

export default VideoP;
