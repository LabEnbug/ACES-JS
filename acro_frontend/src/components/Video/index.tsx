import { useEffect, useRef, useState } from "react";
import videojs from "video.js";
import styles from './style/index.module.less';
import 'video.js/dist/video-js.css';
import { createCanvas, loadImage } from 'canvas'; 
import SideBar from './sidebar';
import FootBar from './footbar';
import BriefIntri from './brief_intro'
import GetAxios from '@/utils/getaxios';
import { Message } from '@arco-design/web-react';
import locale from './locale';
import useLocale from '@/utils/useLocale';
import cs from 'classnames';
import CommentPop from '@/components/Comment'
function VideoPlayer({
  hlsPlayList,
  playIndex,
  reflectPlayIndex,
  recordWatched,
  options,
  ...props
}) {
  const t = useLocale(locale);
  const videoRef = useRef(null);
  const playerRef = useRef(null);
  const [currentVideoIndex, setCurrentVideoIndex] = useState(playIndex);
  const canvas = createCanvas(400, 400);
  const ctx = canvas.getContext('2d');
  const [footBarVis, setFootBarVis] = useState(false);
  const [playstate, setPlayState] = useState(false);
  const [timestate, setTimeState] = useState({
    'now': 0,
    'whole': 0
  });
  const [autoNext, setAutoNext] = useState(true);
  const [volume, setVolume] = useState(0);
  const [playrate, setPlayRate] = useState(1);
  const [fullscreen, setFullScreen] = useState(false);
  const [userfavorite, SetUserfavorite] = useState(false);
  const [userlike, SetUserLike] = useState(false);
  const [likecount, SetLikeCount] = useState(0);
  const [favoritecount, SetFavoriteCount] = useState(0);
  const [forwardedcount, SetForwardedCount] = useState(0);
  const [commentedcount, SetCommentedCount] = useState(0);
  const [follow, SetFollow] = useState(false);
  const [commentvis, SetCommentVis] = useState(false);
  const [backgroundimage, SetBackGroundImage] = useState('');
  const [videoinfo, setVideoInfo ] = useState({
    nickname: 'default',
    username: 'default',
    content: 'default',
    be_watched_count: 0,
    time: "2023-10-31T18:43:57.000Z",
    video_uid:  null,
    keyword: '#default',
    user_id: -1
  });
  let clickTimeout = useRef(null);
  const JudgeStatus = (data: any) => {
    if (data.status != 200) {
      // Message.error(t['message.notfind'])
      return false;
    }
    return true;
  }

  const OpenComments = () => {
    const imag = playerRef.current.el().style.backgroundImage;
    SetCommentVis((pre)=>(!pre));
    SetBackGroundImage(imag);
  }

  const handlePlayerClick = () => {
    // 如果我们已经有一个等待的单击（意味着这可能是一个双击）
    if (clickTimeout.current !== null) {
      clearTimeout(clickTimeout.current); // 清除定时器
      clickTimeout.current = null;
    } else {
      // 如果还没有等待的单击（意味着这是第一次点击）
      clickTimeout.current = setTimeout(() => {
        playerRef.current.paused() ? playerRef.current.play() : playerRef.current.pause();
        clickTimeout.current = null;
      }, 250); // 300ms的延迟来检测是否有第二次点击（双击）
    }
  };

  const getVideoInfo = (uid)=> {
    const baxios = GetAxios();
    const param1 = new FormData();
    param1.append('video_uid',  uid);
    baxios.post('v1-api/v1/video/info', param1).then(res => {
        if (JudgeStatus(res.data)) {
          const video = res.data.data.video;
          setVideoInfo({
            nickname: video['user']['nickname'],
            username: video['user']['username'],
            content: video['content'],
            be_watched_count: video['be_watched_count'],
            video_uid:  video['video_uid'],
            time:  video['upload_time'],
            keyword: video['keyword'],
            user_id: video['user']['user_id'],
          });
          SetUserfavorite(video['is_user_favorite']);
          SetUserLike(video['is_user_liked']);
          SetFavoriteCount(video['be_favorite_count']);
          SetLikeCount(video['be_liked_count']);
          SetForwardedCount(video['be_forwarded_count']);
          SetCommentedCount(video['be_commented_count']);
          SetFollow(video['user']['be_followed']);
          window.localStorage.setItem('is_user_favorite', video['is_user_favorite']);
          window.localStorage.setItem('is_user_like', video['is_user_liked']);
          window.localStorage.setItem('follow',  video['user']['be_followed']);
        }
      }).catch(error => {
        console.error(error)
      });
  }

  const changefollow = ()=> {
    const status = window.localStorage.getItem('follow') == null ? false : JSON.parse(window.localStorage.getItem('follow'));
    const action = status ? `unfollow` : `follow`;
    const param = new FormData();
    const baxios = GetAxios();
    param.append('action',  action);
    param.append('user_id', (videoinfo['user_id']).toString());

    baxios.post('v1-api/v1/user/follow', param).then(res=> {
      if (JudgeStatus(res.data)) {
        window.localStorage.setItem(`follow`, (!status).toString());
        SetFollow(!status);
      } else {
        Message.error(t['message.notlog'])
      }
    }).catch(e => {
      console.error(e);
    })
  }
  
  const clickCount = (a_type, setS, setC)=> {
    const item_name = `is_user_${a_type}`;
    const status = window.localStorage.getItem(item_name) == null ? false : JSON.parse(window.localStorage.getItem(item_name));
    const action = status ? `un${a_type}` : `${a_type}`;
    const param = new FormData();
    const baxios = GetAxios();

    param.append('action',  action);
    param.append('video_uid', videoinfo['video_uid']);
    baxios.post('v1-api/v1/video/action', param).then(res=> {
      if (JudgeStatus(res.data)) {
        if (status) {
          setC((pre)=>(pre-1)) ;
        } else {
          setC((pre)=>(pre+1)) ;
        }
        setS(!status);
        window.localStorage.setItem(`is_user_${a_type}`, (!status).toString());
      } else {
        Message.error(t['message.notlog']);
      }
    }).catch(e => {
      console.error(e);
    })
  };

  const videoDoubleClick = () => {
    const item_name = 'is_user_like';
    const param = new FormData();
    const baxios = GetAxios();
    const status = window.localStorage.getItem(item_name) == null ? false : JSON.parse(window.localStorage.getItem(item_name));
    param.append('action',  'like');
    param.append('video_uid', videoinfo['video_uid']);
    baxios.post('v1-api/v1/video/action', param).then(res=> {
      if (JudgeStatus(res.data)) {
        SetUserLike(true);
        window.localStorage.setItem(item_name, (true).toString());
        status ? true :  SetLikeCount((pre)=> pre+1)
      } else {
        Message.error(t['message.notlog']);
      }
    }).catch(e => {
      console.error(e);
    })
  }

  const clickfoward = ()=> {
    const param = new FormData();
    const baxios = GetAxios();
    param.append('video_uid',  videoinfo['video_uid']);
    const currentURL = window.location.href;
    const textArea = document.createElement("textarea");
    textArea.value = currentURL;
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    try {
      document.execCommand('copy');
      baxios.post('v1-api/v1/video/forward', param).then(res=> {}).catch(e => {
        console.error(e);
      });
      SetForwardedCount((pre)=>(pre+1))
      Message.info(t['message.share']);
    } catch (err) {
      // console.error('Unable to copy to clipboard', err);
      Message.error(t['message.share.failed']);
    }
    document.body.removeChild(textArea);
  }

  const changeFullScreen = () => {
    setFullScreen(!playerRef.current.isFullscreen())
    if (playerRef.current.isFullscreen()) {
        playerRef.current.exitFullscreen();
    } else {
        playerRef.current.requestFullscreen();
    }
  }

  const clickPlay = () => {
    if (!playstate) {
      playerRef.current.pause();
    } else {
      playerRef.current.play();
    }
  }

  const setAuto = (e) => {
    setAutoNext(e);
    window.localStorage.setItem('autonext', e);
  }

  const changeVolume = (e) => {
    playerRef.current.volume(e/100);
  }

  const setPlayBackRate = (e) => {
    playerRef.current.playbackRate(e)
  }

  useEffect(() => {
    const upDateBackGround = (playerRef, url)=> {
      // 创建一个 Image 对象
      const img = new Image();
      // 设置图像的加载完成回调
      img.setAttribute("crossOrigin",'Anonymous')
      img.onload = () => {
        ctx.filter  = 'blur(50px)'; // 例如，应用灰度滤镜
        ctx.drawImage(img, 0, 0, 400, 400);

      // 在图像上应用滤镜效果
      // 将处理后的图像数据作为背景图片
        const filteredImageData = ctx.canvas.toDataURL('image/jpeg');
        // playerRef.current.el().classList.add(containerStyle);
        // 图像加载完成后，将其设置为背景图像
        // playerRef.current.el().style.backgroundColor = 'blue';
        playerRef.current.el().style.backgroundImage = `url(${filteredImageData})`;
        SetBackGroundImage(`url(${filteredImageData})`);
        // playerRef.current.el().style.filter = 'blur(10px)';
      };

      // 设置图像的加载失败回调（可选）
      img.onerror = () => {
        console.error('Failed to load image');
      };

      // 开始加载图像
      img.src = url;
    }
    console.log(hlsPlayList)
    const realindex = currentVideoIndex >= 0 ? currentVideoIndex % hlsPlayList.length : currentVideoIndex % hlsPlayList.length + hlsPlayList.length;
    reflectPlayIndex(Number.isNaN(realindex) ? 0 : realindex);
    if (!playerRef.current && videoRef.current && hlsPlayList.length > 0) {
      playerRef.current = videojs(videoRef.current, {
        crossOrigin: "Anonymous",
        controls: true,
        sources: [{ src: hlsPlayList[realindex]['play_url'], type: "application/x-mpegURL" }],
        poster: hlsPlayList[realindex]['cover_url'],
        preload: "auto",
        autoplay: true,
        userActions: {
          doubleClick: false, // 值也可以是一个函数
          click: false,
        },
        ...options
      });

      playerRef.current.on('ended', () => {
        const autoNext = window.localStorage.getItem('autonext') ? JSON.parse(window.localStorage.getItem('autonext')) : false;
        if (autoNext) {
          setCurrentVideoIndex((prevIndex) => prevIndex + 1);
        } else {
          playerRef.current.currentTime(0);
          playerRef.current.play();
        }
      });
      playerRef.current.on('ready', () => {
        setFullScreen(playerRef.current.isFullscreen());
        setVolume(Math.floor(playerRef.current.volume()*100));
      });
      playerRef.current.on('play', ()=> {
        setFootBarVis(true);
        setPlayState(false);
        recordWatched();
      })
      playerRef.current.on('pause', ()=> {
        setPlayState(true)
      })
      playerRef.current.on('timeupdate', function() {
        const currentPlayTime = playerRef.current.currentTime();
        const totalDuration = playerRef.current.duration();
        setTimeState({
          'now': currentPlayTime,
          'whole': totalDuration
        })
      });
      playerRef.current.on('volumechange', function() {
        setVolume(Math.floor(playerRef.current.volume() * 100))
      });

      playerRef.current.on('ratechange', function() {
        setPlayRate(playerRef.current.playbackRate());
      });

      playerRef.current.on('fullscreenchange', function() {
        setFullScreen(playerRef.current.isFullscreen());
      });

      // playerRef.current.on('click', handlePlayerClick);
      
      playerRef.current.el().classList.add(styles['video-background']);
      playerRef.current.controlBar.getChild('playToggle').hide();
      playerRef.current.controlBar.getChild('VolumePanel').hide();
      playerRef.current.controlBar.getChild('FullscreenToggle').hide();
      playerRef.current.controlBar.getChild('RemainingTimeDisplay').hide();
      playerRef.current.controlBar.removeChild('pictureInPictureToggle');
      upDateBackGround(playerRef, hlsPlayList[realindex]['cover_url'])
    } else if (playerRef.current && videoRef.current && hlsPlayList.length != 0) {
      playerRef.current.src({
        src: hlsPlayList[realindex]['play_url'], type: "application/x-mpegURL"
      })

      playerRef.current.poster(hlsPlayList[realindex]['cover_url'])
      upDateBackGround(playerRef, hlsPlayList[realindex]['cover_url'])
    }
    if (hlsPlayList.length > 0) {
      getVideoInfo(hlsPlayList[realindex]['video_uid'])
      console.log(hlsPlayList[realindex]);
    }
  }, [hlsPlayList, options, currentVideoIndex]);

  useEffect(() => {
    const player = playerRef.current;

    return () => {
      if (player && !player.isDisposed()) {
        console.log(123)
        player.dispose();
        playerRef.current = null;
      }
    };
  }, [playerRef]);

  useEffect(() => {
    const handleMouseWheel = (event) => {
      const specifiedArea = document.getElementById('specified-area');
      if (specifiedArea && specifiedArea.contains(event.target)) {
        if (event.deltaY > 5) {
          setCurrentVideoIndex((prevIndex) =>
            prevIndex + 1
          );
        } else if (event.deltaY < -5) {
          setCurrentVideoIndex((prevIndex) => prevIndex > 0 ? (prevIndex - 1) : 0);
        }
      }
    };
    window.addEventListener('wheel', handleMouseWheel);
    return () => {
      window.removeEventListener('wheel', handleMouseWheel);
    };
  }, [hlsPlayList, currentVideoIndex]);

  useEffect(() => {
    const handleKeyDown = (event) => {
      const specifiedArea = document.getElementById('specified-area');
      if (specifiedArea && specifiedArea.contains(event.target)) {
        if (event.keyCode == 32) {
          if (!playstate) {
            playerRef.current.pause();
          } else {
            playerRef.current.play();
          }
        }
      }
      if (event.key === 'ArrowUp') {
        setCurrentVideoIndex((prevIndex) => prevIndex > 0 ? prevIndex - 1 : 0);
      } else if (event.key === 'ArrowDown') {
        setCurrentVideoIndex((prevIndex) => prevIndex + 1);
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [playstate, currentVideoIndex, playerRef]);



  return (
    <div style={{display: 'inline-flex', width: '100%', height: '100%'}}>
      <div data-vjs-player className={styles['video-container']} >
        <video ref={videoRef} onClick={handlePlayerClick} onDoubleClick={videoDoubleClick} id="specified-area" className={`vjs-default-skin video-js ${styles['video-pos-js-9-16']}`} controls></video>
        <SideBar videoinfo={videoinfo} 
                 userfavorite={userfavorite} 
                 userlike={userlike} 
                 ikecount={likecount} 
                 favoritecount={favoritecount}  
                 forwardedcount={forwardedcount}
                 commentedcount={commentedcount}
                 clickfavorite={{'func': clickCount, 'params': ['favorite', SetUserfavorite, SetFavoriteCount]}}
                 clicklike={{'func': clickCount, 'params': ['like', SetUserLike, SetLikeCount]}}
                 clickfoward={clickfoward}
                 followed={follow}
                 changefollow={changefollow}
                 cilckcomment={OpenComments}
        />
        <FootBar id='footbar'
                 ref={playerRef}
                 visible={footBarVis}
                 playstate={playstate}
                 timestate={timestate}
                 playclick={clickPlay}
                 volume={volume}
                 volumechange={changeVolume}
                 setauto={setAuto}
                 autostate={autoNext}
                 playbackrate={playrate}
                 setplaybackrate={setPlayBackRate}
                 fullscreen={fullscreen}
                 fullscreenchange={changeFullScreen}
               />
        <BriefIntri videoinfo={videoinfo} />
      </div>
      <div className={commentvis ? cs(styles['comment-container-vis']) : styles['comment-container-dis']} style={commentvis ? {backgroundImage: backgroundimage} : {}}>
        <CommentPop videoinfo={videoinfo} />
      </div>
    </div>
  );
};

export default VideoPlayer;