import { useEffect, useRef, useState } from "react";
import videojs from "video.js";
import styles from './style/index.module.less';
import 'video.js/dist/video-js.css';
import { createCanvas, loadImage } from 'canvas'; 
import SideBar from './sidebar';
import FootBar from './footbar';

function VideoPlayer({
  hlsPlayList,
  playindex,
  options,
  ...props
}) {
  const videoRef = useRef(null);
  const playerRef = useRef(null);
  const [currentVideoIndex, setCurrentVideoIndex] = useState(playindex);
  const canvas = createCanvas(400, 400);
  const ctx = canvas.getContext('2d');
  const [footbarVis, setFootBarVis] = useState(false);
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
    if (!playerRef.current && videoRef.current && hlsPlayList.length > playindex) {
      playerRef.current = videojs(videoRef.current, {
        crossOrigin: "Anonymous",
        controls: true,
        sources: [{ src: hlsPlayList[currentVideoIndex]['play_url'], type: "application/x-mpegURL" }],
        poster: hlsPlayList[currentVideoIndex]['cover_url'],
        preload: "auto",
        autoplay: false,
        ...options
      });
  
      playerRef.current.on('ended', () => {
        // if (currentVideoIndex < hlsPlayList.length - 1) {
        //   setCurrentVideoIndex((prevIndex) => prevIndex + 1);
        // } else {
        //   setCurrentVideoIndex(0);
        // }
        playerRef.current.currentTime = 0; 
        playerRef.current.play();
      });
      playerRef.current.el().classList.add(styles['video-background']); 
      playerRef.current.controlBar.getChild('playToggle').hide();
      playerRef.current.controlBar.getChild('VolumePanel').hide();
      playerRef.current.controlBar.getChild('FullscreenToggle').hide();
      playerRef.current.controlBar.getChild('RemainingTimeDisplay').hide();
      playerRef.current.controlBar.removeChild('pictureInPictureToggle');
      upDateBackGround(playerRef, hlsPlayList[currentVideoIndex]['cover_url'])
    } else if (playerRef.current && videoRef.current && hlsPlayList.length > playindex) {
      console.log(currentVideoIndex)
      playerRef.current.src({
        src: hlsPlayList[currentVideoIndex]['play_url'], type: "application/x-mpegURL" 
      })
      playerRef.current.poster(hlsPlayList[currentVideoIndex]['cover_url'])
      upDateBackGround(playerRef, hlsPlayList[currentVideoIndex]['cover_url'])
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
            prevIndex < hlsPlayList.length - 1 ? prevIndex + 1 : prevIndex
          );
        } else if (event.deltaY < -5) {
          setCurrentVideoIndex((prevIndex) => (prevIndex > 0 ? prevIndex - 1 : prevIndex));
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
      if (event.key === 'ArrowUp' && currentVideoIndex > 0) {
        setCurrentVideoIndex((prevIndex) => prevIndex - 1);
      } else if (event.key === 'ArrowDown' && currentVideoIndex < hlsPlayList.length - 1) {
        setCurrentVideoIndex((prevIndex) => prevIndex + 1);
      }
    };
    window.addEventListener('keydown', handleKeyDown);

    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [hlsPlayList, currentVideoIndex]);

  return (
    <div data-vjs-player className={styles['video-container']} >
       <video ref={videoRef} id="specified-area" className={`vjs-default-skin ${styles['video-pos-js-9-16']} video-js`} controls></video>
       <SideBar/>
       <FootBar playRef={playerRef} visible={footbarVis}/>
    </div>
  );
};

export default VideoPlayer;