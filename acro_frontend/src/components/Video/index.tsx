import { useEffect, useRef, useState } from "react";
import videojs from "video.js";
import styles from './style/index.module.less';
import 'video.js/dist/video-js.css';
import { createCanvas, loadImage } from 'canvas'; 
import SideBar from './sidebar';
import FootBar from './footbar';
import BriefIntri from './brief_intro'

function VideoPlayer({
  hlsPlayList,
  playIndex,
  reflectPlayIndex,
  recordWatched,
  options,
  ...props
}) {
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
  const [videoinfo, setVideoInfo ] = useState({
    nickname: 'default',
    username: 'default',
    content: 'default',
    be_commented_count: 0,
    be_favorite_count: 0,
    be_liked_count: 0,
    be_forwarded_count: 0,
    be_watched_count: 0,
    time: "2023-10-31T18:43:57.000Z",
    video_uid:  null,
    keyword: '#default',
  });


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
        playerRef.current.el().style.backgroundColor = 'blue';
        playerRef.current.el().style.backgroundImage = `url(${filteredImageData})`;
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
    reflectPlayIndex(realindex);
    if (!playerRef.current && videoRef.current && hlsPlayList.length > 0) {
      playerRef.current = videojs(videoRef.current, {
        crossOrigin: "Anonymous",
        controls: true,
        sources: [{ src: hlsPlayList[realindex]['play_url'], type: "application/x-mpegURL" }],
        poster: hlsPlayList[realindex]['cover_url'],
        preload: "auto",
        autoplay: true,
        ...options
      });

      playerRef.current.on('ended', () => {
        // if (currentVideoIndex < hlsPlayList.length - 1) {
        //   setCurrentVideoIndex((prevIndex) => prevIndex + 1);
        // } else {
        //   setCurrentVideoIndex(0);
        // }
        const autoNext = window.localStorage.getItem('autonext') ? JSON.parse(window.localStorage.getItem('autonext')) : false;
        if (autoNext) {
          setCurrentVideoIndex((prevIndex) => prevIndex + 1);
          // console.log(21371983791)
          // playerRef.current.play();
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
      console.log(hlsPlayList[realindex])
      setVideoInfo({
        nickname: hlsPlayList[realindex]['user']['nickname'],
        username: hlsPlayList[realindex]['user']['username'],
        content: hlsPlayList[realindex]['content'],
        be_commented_count: hlsPlayList[realindex]['be_commented_count'],
        be_favorite_count: hlsPlayList[realindex]['be_favorite_count'],
        be_liked_count: hlsPlayList[realindex]['be_liked_count'],
        be_forwarded_count: hlsPlayList[realindex]['be_forwarded_count'],
        be_watched_count: hlsPlayList[realindex]['be_watched_count'],
        video_uid:  hlsPlayList[realindex]['video_uid'],
        time:  hlsPlayList[realindex]['upload_time'],
        keyword: hlsPlayList[realindex]['keyword'],
      })
      console.log(hlsPlayList[realindex]);
    }
  }, [hlsPlayList, options, currentVideoIndex]);

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
          setCurrentVideoIndex((prevIndex) =>
            prevIndex + 1
          );
        } else if (event.deltaY < -5) {
          setCurrentVideoIndex((prevIndex) => (prevIndex - 1));
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
      if (event.key === 'ArrowUp') {
        setCurrentVideoIndex((prevIndex) => prevIndex - 1);
      } else if (event.key === 'ArrowDown') {
        setCurrentVideoIndex((prevIndex) => prevIndex + 1);
      }
    };
    window.addEventListener('keydown', handleKeyDown);

    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [hlsPlayList, currentVideoIndex]);

  return (
    <>
      <div data-vjs-player className={styles['video-container']} >
        <video ref={videoRef} id="specified-area" className={`vjs-default-skin video-js ${styles['video-pos-js-9-16']}`} controls></video>
        <SideBar videoinfo={videoinfo} />
        <FootBar id='footbar'
                 ref={playerRef}
                 visible={footBarVis}
                 playstate={playstate}
                 timestate={timestate}
                 playclick={clickPlay}
                 volume={volume}
                 volumechange={changeVolume}
                 setauto={setAuto}
                 autoNext={autoNext}
                 playbackrate={playrate}
                 setplaybackrate={setPlayBackRate}
                 fullscreen={fullscreen}
                 fullscreenchange={changeFullScreen}
               />
        <BriefIntri videoinfo={videoinfo} />
      </div>
    </>
  );
};

export default VideoPlayer;