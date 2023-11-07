import React, {useRef} from 'react';
import { Message } from '@arco-design/web-react';
import locale from './locale';
import useLocale from '@/utils/useLocale';
import styles from './style/index.module.less';
import VideoPlayer from '@/components/Video';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import GetVideType from '@/utils/getvideotype';
import baxios from "@/utils/getaxios";
import Head from "next/head";

interface VideoPlayerMethods {
  handleKeyDown: (event: KeyboardEvent) => void;
}
function VideoP() {
  const default_type = 'comprehensive';
  const t = useLocale(locale);
  const tg = useLocale();
  const router = useRouter();
  const videoPlayerRef = useRef<VideoPlayerMethods>(null);
  const type = router.query['type'] || default_type;
  const video_uid = router.query['video_uid'];
  const rel = router.query['rel'] || '';
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
    const pre_type = window.sessionStorage.getItem('pretype') || 'comprehensive';
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
      window.sessionStorage.setItem('playvideo-id', playlist[index]['video_uid']);
    }
  };

  const recordWatched = () => {
    const pre = window.sessionStorage.getItem('playvideo-pre-id');
    const uid = window.sessionStorage.getItem('playvideo-id');
    if (uid && pre != uid) {
      baxios
        .post('/videos/' + uid.toString() + '/actions/' + 'watch')
        .then((response) => {
          window.sessionStorage.setItem('playvideo-pre-id', uid);
        })
        .catch((error) => {
          console.error(error);
        });
    }
  };

  useEffect(() => {
    const pre_type = window.sessionStorage.getItem('pretype');
    // console.log(type, video_uid, playIndex);
    if (type != default_type && GetVideType(type) >= 999) {
      router.push('/video');
    }

    if (rel && rel === '1') {
      baxios
        .get('/videos/' + video_uid.toString())
        .then((response1) => {
          if (JudgeStatus(response1.data)) {
            // add this video to next play in playlist
            const tmp_list = [...playlist];
            tmp_list.splice(playIndex + 1, 0, response1.data.data.video);
            SetPlayList(tmp_list);
            // SetPlayIndex(playIndex + 1);
            console.log(tmp_list);
            // press keydown to play next video
            if (videoPlayerRef.current) {
              const fakeEvent = new KeyboardEvent("keydown", {
                key: "ArrowDown",
                code: "ArrowDown",
                keyCode: 40,
                charCode: 0,
                view: window
              });
              videoPlayerRef.current.handleKeyDown(fakeEvent);
            }
          }
        })
        .catch((error) => {
          console.error(error);
        });
    } else if (playlist.length == 0 || pre_type != type) {
      // param.append('page', page)
      baxios
        .get('/videos?' + 'limit=' + limit +
          (type != default_type ? '&type=' + GetVideType(type) : ''))
        .then((response) => {
          const data = response.data;
          window.sessionStorage.setItem('pretype', type.toString());
          if (JudgeStatus(data)) {
            if (
              video_uid &&
              data.data.video_list.length > 0 &&
              video_uid != data.data.video_list[0]['video_uid']
            ) {
              baxios
                .get('/videos/' + video_uid.toString())
                .then((response1) => {
                  if (JudgeStatus(response1.data)) {
                    data.data.video_list.unshift(response1.data.data.video);
                  }
                  SetPlayList(data.data.video_list);
                  if (videoPlayerRef.current) {
                    const fakeEvent = new KeyboardEvent("keydown", {
                      key: "ArrowDown",
                      code: "ArrowDown",
                      keyCode: 40,
                      charCode: 0,
                      view: window
                    });
                    videoPlayerRef.current.handleKeyDown(fakeEvent);
                  }
                })
                .catch((error) => {
                  console.error(error);
                });
              //
            } else {
              SetPlayIndex(0);
              SetPlayList(data.data.video_list);
              if (videoPlayerRef.current) {
                const fakeEvent = new KeyboardEvent("keydown", {
                  key: "ArrowDown",
                  code: "ArrowDown",
                  keyCode: 40,
                  charCode: 0,
                  view: window
                });
                videoPlayerRef.current.handleKeyDown(fakeEvent);
              }
            }
          }
        })
        .catch((error) => {
          console.error(error);
        });
    } else if (playIndex >= playlist.length - 3) {
      baxios
        .get('/videos?' + 'limit=' + limit + '&' +
          (type != default_type ? 'type=' + GetVideType(type) + '&' : '') +
          'start=' + playlist.length.toString())
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
  }, [type, video_uid, playIndex]);

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
    <>
      <Head>
        <title>{tg['title.global']}</title>
      </Head>
      <div className={styles.container}>
        <div className={styles.wrapper}>
          <VideoPlayer
            ref={videoPlayerRef}
            hlsPlayList={playlist}
            playIndex={playIndex}
            reflectPlayIndex={reflectPlayIndex}
            recordWatched={recordWatched}
            // options={undefined}
          />
        </div>
      </div>
    </>
  );
}

export default VideoP;
