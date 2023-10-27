"use client"

import { FC, useEffect, useRef, useState } from "react";
import { Box, BoxProps } from "@mantine/core";
import videojs from "video.js";
import Player from "video.js/dist/types/player";
import "video.js/dist/video-js.css";
import "@/app/css/video.css"
import type { StaticImageData } from 'next/image'

interface VideoPlayerProps extends BoxProps {
  hlsPlayList: any;
  playindex: any;
  options?: Player["options"];
}

export const VideoPlayer: FC<VideoPlayerProps> = ({
  hlsPlayList,
  playindex,
  options,
  ...props
}) => {
  const videoRef = useRef<HTMLDivElement>(null);
  const playerRef = useRef<Player | null>(null);
  const [currentVideoIndex, setCurrentVideoIndex] = useState(playindex);

  useEffect(() => {
    if (!playerRef.current && videoRef.current && hlsPlayList.length > playindex) {
      const videoElement = document.createElement("video-js");

      videoElement.classList.add("vjs-default-skin");
      videoElement.classList.add("video-pos-js-9-16");
      
      
      videoRef.current.appendChild(videoElement);

      
      playerRef.current = videojs(videoElement, {
        controls: true,
        sources: [{ src: hlsPlayList[currentVideoIndex]['play_url'], type: "application/x-mpegURL" }],
        poster: hlsPlayList[currentVideoIndex]['cover_url'],
        preload: "auto",
        autoplay: true,
        ...options
      });
    }
    if (playerRef.current && videoRef.current) {
      console.log(currentVideoIndex)
      playerRef.current.src({
        src: hlsPlayList[currentVideoIndex]['play_url'], type: "application/x-mpegURL" 
      })
      playerRef.current.poster(hlsPlayList[currentVideoIndex]['cover_url'])
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
        if (event.deltaY > 0) {
          setCurrentVideoIndex((prevIndex) =>
            prevIndex < hlsPlayList.length - 1 ? prevIndex + 1 : prevIndex
          );
        } else if (event.deltaY < 0) {
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
    <div id="specified-area" className="video-container justify-center">
      <Box data-vjs-player {...props}  >
        <Box ref={videoRef} />
      </Box>
    </div>
  );
};
