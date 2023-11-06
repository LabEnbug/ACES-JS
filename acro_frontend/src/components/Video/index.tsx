import { useEffect, useRef, useState } from 'react';
import videojs from 'video.js';
import styles from './style/index.module.less';
import 'video.js/dist/video-js.css';
import { createCanvas, loadImage } from 'canvas';
import SideBar from './sidebar';
import FootBar from './footbar';
import BriefIntri from './brief_intro';
import { Message, Tooltip } from '@arco-design/web-react';
import { Like} from '@icon-park/react'
import locale from './locale';
import useLocale from '@/utils/useLocale';
import cs from 'classnames';
import SiderTabs from '@/components/SiderTabs';
import baxios from "@/utils/getaxios";
import BulletScreen, { StyledBullet } from 'rc-bullets';

function VideoPlayer({
  hlsPlayList,
  playIndex,
  reflectPlayIndex,
  recordWatched,
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
    now: 0,
    whole: 0,
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
  const screen = useRef(null);
  const bullets = useRef({index:0, all:[]});
  const [showbullet, setShowBullet] = useState(false);
  const [hearts, setHearts] = useState([]);

  const getRandomMargin = () => {
    return (bullets.current.index % 4) * 40;
  }

  const [videoinfo, setVideoInfo] = useState({
    nickname: 'default',
    username: 'default',
    content: 'default',
    be_watched_count: 0,
    time: '2023-10-31T18:43:57.000Z',
    video_uid: null,
    keyword: '#default',
    user_id: -1,
  });
  const clickTimeout = useRef(null);
  const clickTimeoutBullet = useRef(null);
  const closeBullet = ()=>{
    console.log(screen.current)
    if (!screen.current) {
      return;
    } 
    const state = screen.current.allHide;
    if (state) {
      screen.current.show();
    } else {
      screen.current.hide();
    }
    setShowBullet(screen.current.allHide);
  };

  const JudgeStatus = (data: any) => {
    if (data.status != 200) {
      // Message.error(t['message.notfind'])
      return false;
    }
    return true;
  };

  const OpenComments = () => {
    const imag = playerRef.current.el().style.backgroundImage;
    SetCommentVis((pre) => !pre);
    SetBackGroundImage(imag);
  };

  const handlePlayerClick = () => {
    // 如果我们已经有一个等待的单击（意味着这可能是一个双击）
    if (clickTimeout.current !== null) {
      clearTimeout(clickTimeout.current); // 清除定时器
      clickTimeout.current = null;
    } else {
      // 如果还没有等待的单击（意味着这是第一次点击）
      clickTimeout.current = setTimeout(() => {
        playerRef.current.paused()
          ? playerRef.current.play()
          : playerRef.current.pause();
        clickTimeout.current = null;
      }, 250); // 300ms的延迟来检测是否有第二次点击（双击）
    }
  };

  const getVideoInfo = (uid) => {
    baxios
      .get('/v1-api/v1/videos/' + uid.toString())
      .then((res) => {
        if (JudgeStatus(res.data)) {
          const video = res.data.data.video;
          setVideoInfo({
            nickname: video['user']['nickname'],
            username: video['user']['username'],
            content: video['content'],
            be_watched_count: video['be_watched_count'],
            video_uid: video['video_uid'],
            time: video['upload_time'],
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
          window.localStorage.setItem(
            'is_user_favorite',
            video['is_user_favorite']
          );
          window.localStorage.setItem('is_user_like', video['is_user_liked']);
          window.localStorage.setItem('follow', video['user']['be_followed']);
        }
      })
      .catch((error) => {
        console.error(error);
      });
  };

  const changefollow = () => {
    const status =
      window.localStorage.getItem('follow') == null
        ? false
        : JSON.parse(window.localStorage.getItem('follow'));

    (status ? baxios.delete : baxios.post)
    ('/v1-api/v1/users/' + videoinfo.username + '/follow')
      .then((res) => {
        if (JudgeStatus(res.data)) {
          window.localStorage.setItem(`follow`, (!status).toString());
          SetFollow(!status);
        } else {
          Message.error(t['message.notlog']);
        }
      })
      .catch((e) => {
        console.error(e);
      });
  };

  const clickCount = (a_type, setS, setC) => {
    const item_name = `is_user_${a_type}`;
    const status =
      window.localStorage.getItem(item_name) == null
        ? false
        : JSON.parse(window.localStorage.getItem(item_name));
    (status ? baxios.delete : baxios.post)
    ('v1-api/v1/videos/' + videoinfo['video_uid'] + "/actions/" + a_type)
      .then((res) => {
        if (JudgeStatus(res.data)) {
          if (status) {
            setC((pre) => pre - 1);
          } else {
            setC((pre) => pre + 1);
          }
          setS(!status);
          window.localStorage.setItem(
            `is_user_${a_type}`,
            (!status).toString()
          );
        } else {
          Message.error(t['message.notlog']);
        }
      })
      .catch((e) => {
        console.error(e);
      });
  };

  const videoDoubleClick = (e) => {
    const item_name = 'is_user_like';
    const status =
      window.localStorage.getItem(item_name) == null
        ? false
        : JSON.parse(window.localStorage.getItem(item_name));
    baxios
      .post('v1-api/v1/videos/' + videoinfo['video_uid'] + '/actions/' + 'like')
      .then((res) => {
        if (JudgeStatus(res.data)) {
          SetUserLike(true);
          window.localStorage.setItem(item_name, true.toString());
          status ? true : SetLikeCount((pre) => pre + 1);
        } else {
          Message.error(t['message.notlog']);
        }
      })
      .catch((e) => {
        console.error(e);
      });
      console.log(e);
      const x = e.nativeEvent.layerX;
      const y = e.nativeEvent.layerY;
      // 创建一个新的心形并设置位置
      const id1 = Math.random();
      const id2 = Math.random();
      const newHeart = { id:  id1, ele: (<Like id={id1.toString()} theme="filled" size="36" fill="red"/>), style: { left: x, top: y, position: 'absolute' }};
      // 添加新的心形到数组中，并在一段时间后移除
      setHearts([...hearts, newHeart]);
      setTimeout(() => {
        setHearts(pre=>([...pre, { id:  id2, ele: (<Like id={id2.toString()} theme="filled" size="36" fill="red"/>), style: { left: x, top: y, position: 'absolute' }}]));
      } , 100);
      setTimeout(() => {
        setHearts((currentHearts) => currentHearts.filter(heart => ![id1, id2].includes(heart.id)));
      }, 1500); // 动画持续时间后移除爱心
  };

  const clickfoward = () => {
    const currentURL = window.location.href;
    const textArea = document.createElement('textarea');
    textArea.value = currentURL;
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    try {
      document.execCommand('copy');
      baxios
        .post('v1-api/v1/videos/' + videoinfo['video_uid'] + '/actions/' + 'forward')
        .then((res) => {})
        .catch((e) => {
          console.error(e);
        });
      SetForwardedCount((pre) => pre + 1);
      Message.info(t['message.share']);
    } catch (err) {
      // console.error('Unable to copy to clipboard', err);
      Message.error(t['message.share.failed']);
    }
    document.body.removeChild(textArea);
  };

  const changeFullScreen = () => {
    setFullScreen(!playerRef.current.isFullscreen());
    if (playerRef.current.isFullscreen()) {
      playerRef.current.exitFullscreen();
    } else {
      playerRef.current.requestFullscreen();
    }
  };

  const clickPlay = () => {
    const state = playerRef.current.paused();
    if (state) {
      playerRef.current.play();
    } else {
      playerRef.current.pause();
    }
  };

  const setAuto = (e) => {
    setAutoNext(e);
    window.localStorage.setItem('autonext', e);
  };

  const changeVolume = (e) => {
    playerRef.current.volume(e / 100);
  };

  const setPlayBackRate = (e) => {
    playerRef.current.playbackRate(e);
  };

  const generateBullet = (content, isSelf, nickname) => {
    const marginT = getRandomMargin();
    return (
      <Tooltip position='top' trigger='hover' content={ isSelf ? t['tooltip.bullets.me'] : nickname }>
        <div style={{marginTop: `${marginT}px`, background: 'transparent'}} className={isSelf ? styles['bullets-text-container-self'] :  styles['bullets-text-container-other']}>
          <span className={styles['bullets-text-style']} > {content} </span>
        </div>
      </Tooltip>
    )
  }

  const sendBullet = (e) => {
    return new Promise((resolve, reject) => {
      const bullet = e.target.value;

      const param = new FormData();
      param.append('content', bullet);
      param.append('comment_at', playerRef.current.currentTime().toString());
      baxios.post('v1-api/v1/videos/' + videoinfo.video_uid + '/bullet_comments', param).then(res=> {
        if (res.data.status == 200) {
            bullets.current.all.splice(screen.current.bullets.length, 0, res.data.data.bullet_comment);
            screen.current.push(
              generateBullet(bullet, true, '')
            )
          resolve('success');
          return;
        }
        Message.error('无法发送弹幕');
        reject(e);
      }).catch(e=>{
        console.error(e);
        reject(e);
        Message.error('无法发送弹幕');
      })
    });
  }

  useEffect(()=> {
    const registerScreen = ()=> {
      const area = document.getElementById('video-player-container');
      if (area) {
        const s = new BulletScreen(area, {duration:10, top: '10px', loopCount: 1});
        screen.current = s;
        setShowBullet(screen.current.allHide);
      } else {
        setTimeout(registerScreen, 50);
      }
    };
    registerScreen();
  }, []);

  useEffect(()=>{
    if (playerRef.current && screen.current) {
      playerRef.current.on('play', () => {
        screen.current.resume();
        setFootBarVis(true);
        setPlayState(false);
        recordWatched();
      });
      playerRef.current.on('pause', () => {
        // screen.current.pause();
        setPlayState(true);
      });
      playerRef.current.on('timeupdate', function () {
        const currentPlayTime = playerRef.current.currentTime();
        const totalDuration = playerRef.current.duration();
        setTimeState({
          now: currentPlayTime,
          whole: totalDuration,
        });
        if (bullets.current.all.length <= bullets.current.index) 
          return;
        const bull = bullets.current.all[bullets.current.index];
        if (bull.comment_at < currentPlayTime) {
          // console.log(bull);
          // console.log(bullets.current);
          // console.log(screen.current.bullets);
          // console.log(screen.current.bullets.length);
          // console.log(currentPlayTime);
          screen.current.push(generateBullet(bull['content'], bull['user']['is_self'], bull['user']['nickname']));
          bullets.current.index += 1;
        }
      });
    }
  }, [screen.current, playerRef.current])

  useEffect(() => {
    const upDateBackGround = (playerRef, url) => {
      // 创建一个 Image 对象
      const img = new Image();
      // 设置图像的加载完成回调
      img.setAttribute('crossOrigin', 'Anonymous');
      img.onload = () => {
        ctx.filter = 'blur(50px)'; // 例如，应用灰度滤镜
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
    };
    console.log(hlsPlayList);
    const realindex =
      currentVideoIndex >= 0
        ? currentVideoIndex % hlsPlayList.length
        : (currentVideoIndex % hlsPlayList.length) + hlsPlayList.length;
    reflectPlayIndex(Number.isNaN(realindex) ? 0 : realindex);
    if (!playerRef.current && videoRef.current && hlsPlayList.length > 0) {
      playerRef.current = videojs(videoRef.current, {
        crossOrigin: 'Anonymous',
        controls: true,
        sources: [
          {
            src: hlsPlayList[realindex]['play_url'],
            type: 'application/x-mpegURL',
          },
        ],
        poster: hlsPlayList[realindex]['cover_url'],
        preload: 'auto',
        autoplay: true,
        userActions: {
          doubleClick: false, // 值也可以是一个函数
          click: false,
        },
      });
      playerRef.current.on('ended', () => {
        const autoNext = window.localStorage.getItem('autonext')
          ? JSON.parse(window.localStorage.getItem('autonext'))
          : false;
        if (autoNext) {
          setCurrentVideoIndex((prevIndex) => prevIndex + 1);
        } else {
          playerRef.current.currentTime(0);
          playerRef.current.play();
        }
        screen.current.clear();
        bullets.current.index = 0;
      });
      playerRef.current.on('ready', () => {
        setFullScreen(playerRef.current.isFullscreen());
        setVolume(Math.floor(playerRef.current.volume() * 100));
      });
      playerRef.current.on('volumechange', function () {
        setVolume(Math.floor(playerRef.current.volume() * 100));
      });

      playerRef.current.on('ratechange', function () {
        setPlayRate(playerRef.current.playbackRate());
      });

      playerRef.current.on('fullscreenchange', function () {
        setFullScreen(playerRef.current.isFullscreen());
      });

      // playerRef.current.on('click', handlePlayerClick);

      playerRef.current.el().classList.add(styles['video-background']);
      playerRef.current.controlBar.getChild('playToggle').hide();
      playerRef.current.controlBar.getChild('VolumePanel').hide();
      playerRef.current.controlBar.getChild('FullscreenToggle').hide();
      playerRef.current.controlBar.getChild('RemainingTimeDisplay').hide();
      playerRef.current.controlBar.removeChild('pictureInPictureToggle');
      upDateBackGround(playerRef, hlsPlayList[realindex]['cover_url']);
    } else if (
      playerRef.current &&
      videoRef.current &&
      hlsPlayList.length != 0
    ) {
      playerRef.current.src({
        src: hlsPlayList[realindex]['play_url'],
        type: 'application/x-mpegURL',
      });

      playerRef.current.poster(hlsPlayList[realindex]['cover_url']);
      upDateBackGround(playerRef, hlsPlayList[realindex]['cover_url']);
    }
    if (hlsPlayList.length > 0) {
      getVideoInfo(hlsPlayList[realindex]['video_uid']);
    }
  }, [hlsPlayList, currentVideoIndex]);

  useEffect(() => {
    const player = playerRef.current;

    return () => {
      if (player && !player.isDisposed()) {
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
          setCurrentVideoIndex((prevIndex) => prevIndex + 1);
        } else if (event.deltaY < -5) {
          setCurrentVideoIndex((prevIndex) =>
            prevIndex > 0 ? prevIndex - 1 : 0
          );
        }
      }
    };
    window.addEventListener('wheel', handleMouseWheel);
    return () => {
      window.removeEventListener('wheel', handleMouseWheel);
    };
  }, []);

  useEffect(()=>{
    if (!videoinfo.video_uid) return;
    if (!screen.current) return;
    screen.current.clear();
    bullets.current.all.length = 0;
    bullets.current.index = 0;

    if (clickTimeoutBullet.current !== null) {
      clearTimeout(clickTimeoutBullet.current); // 清除定时器
      clickTimeoutBullet.current = null;
    } 
    const fetchMore= (offset)=>{
      clearTimeout(clickTimeoutBullet.current); 
      clickTimeoutBullet.current = null;
      baxios.get(
        'v1-api/v1/videos/' + videoinfo.video_uid + '/bullet_comments' + '?' +
        'limit' + '=' + '50' + '&' +
        'start' + '=' + `${offset}`
      ).then(res=> {
        if (res.data.status == 200) {
          const data = res.data.data;
          if (!data.bullet_comment_list || data.bullet_comment_list.length > 0) {
            clearTimeout(clickTimeout.current); // 清除定时器
            clickTimeout.current = null;
            return;
          }
          if (data.bullet_comment_list) {
            bullets.current.all.push(...data.bullet_comment_list);
          }
          clickTimeout.current = setTimeout(() => {
            fetchMore(bullets.current.all.length);
          }, 500); 
        } else Message.error('Can not fetch bullets');
      }).catch(e=>{
        console.error(e);
        Message.error('Can not fetch bullets');
      });
    }
    fetchMore(0);
  }, [videoinfo.video_uid, screen.current])

  useEffect(() => {
    const handleKeyDown = (event) => {
      const specifiedArea = document.getElementById('specified-area');
      if (specifiedArea && specifiedArea.contains(event.target)) {
        if (event.keyCode == 32) {
          const state = playerRef.current.paused();
          if (!state) {
            playerRef.current.pause();
          } else {
            playerRef.current.play();
          }
        }
      }
      if (event.key === 'ArrowUp') {
        setCurrentVideoIndex((prevIndex) =>
          prevIndex > 0 ? prevIndex - 1 : 0
        );
      } else if (event.key === 'ArrowDown') {
        setCurrentVideoIndex((prevIndex) => prevIndex + 1);
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [ playerRef ]);

  return (
    <div
      style={{
        display: 'inline-flex',
        width: '100%',
        height: '100%',
        minWidth: 680,
      }}
    >
      <div data-vjs-player id='video-player-container' className={styles['video-container']}>
        <video
          ref={videoRef}
          onClick={handlePlayerClick}
          onDoubleClick={videoDoubleClick}
          id="specified-area"
          className={`vjs-default-skin video-js ${styles['video-pos-js-9-16']}`}
          controls
        >
        </video>
        {hearts.map((heart) => (
            <div key={heart.id} className={styles['heart-animate']} style={heart.style}>
              {heart.ele}
            </div>
          ))}
        <SideBar
          videoinfo={videoinfo}
          userfavorite={userfavorite}
          userlike={userlike}
          ikecount={likecount}
          favoritecount={favoritecount}
          forwardedcount={forwardedcount}
          commentedcount={commentedcount}
          clickfavorite={{
            func: clickCount,
            params: ['favorite', SetUserfavorite, SetFavoriteCount],
          }}
          clicklike={{
            func: clickCount,
            params: ['like', SetUserLike, SetLikeCount],
          }}
          clickfoward={clickfoward}
          followed={follow}
          changefollow={changefollow}
          cilckcomment={OpenComments}
        />
        <FootBar
          id="footbar"
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
          video_place='specified-area'
          sendBullet={sendBullet}
          closeBullet={closeBullet}
          bulletState={showbullet}
        />
        <BriefIntri videoinfo={videoinfo} />
      </div>
      <div
        className={
          commentvis
            ? cs(styles['comment-container-vis'])
            : styles['comment-container-dis']
        }
        style={commentvis ? { backgroundImage: backgroundimage } : {}}
      >
        <SiderTabs videoinfo={videoinfo} />
      </div>
    </div>
  );
}

export default VideoPlayer;
