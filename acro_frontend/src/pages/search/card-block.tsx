import React, {useEffect, useRef, useState} from 'react';
import cs from 'classnames';
import {
  Button,
  Switch,
  Tag,
  Card,
  Descriptions,
  Typography,
  Dropdown,
  Menu,
  Skeleton, Avatar,
} from '@arco-design/web-react';
import useLocale from '@/utils/useLocale';
import locale from './locale';
import {VideoCard, UserCard} from './interface';
import styles from './style/index.module.less';
import {useRouter} from "next/router";
import videojs from "video.js";
import {Like} from "@icon-park/react";
import IconButton from "@/components/NavBar/IconButton";
import {IconClockCircle, IconShake} from "@arco-design/web-react/icon";

interface CardBlockType {
  type: 'video' | 'user' ;
  card: VideoCard | UserCard;
  loading?: boolean;
}



function CardBlock(props: CardBlockType) {
  const { type, card = {} } = props;
  const [visible, setVisible] = useState(false);
  const [status, setStatus] = useState(card.status);
  const [loading, setLoading] = useState(props.loading);

  const t = useLocale(locale);

  const router = useRouter();


  function makeNewSearch(keyword: string) {
    // change search input in navbar
    const searchInput = document.getElementById('search-input');
    if (searchInput) {
      searchInput.setAttribute('value', keyword);
    }
    router.push({
      pathname: '/search',
      query: {
        q: keyword,
      },
    });
  }

  function goToVideoPage(video_uid: string) {
    router.push({
      pathname: `/video`,
      query: {
        video_uid: video_uid,
      },
    });
  }


  useEffect(() => {
    setLoading(props.loading);
  }, [props.loading]);


  // parse time "2023-10-27T16:43:57+08:00" string to some like "3 days ago"
  const parseTime = (time: string) => {
    const date = new Date(time);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const days = Math.floor(diff / (24 * 3600 * 1000));
    const hours = Math.floor((diff % (24 * 3600 * 1000)) / (3600 * 1000));
    const minutes = Math.floor((diff % (3600 * 1000)) / (60 * 1000));
    const seconds = Math.floor((diff % (60 * 1000)) / 1000);
    if (days > 0) {
      return `${days} 天前`;
    } else if (hours > 0) {
      return `${hours} 小时前`;
    } else if (minutes > 0) {
      return `${minutes} 分钟前`;
    } else {
      return `${seconds} 秒前`;
    }
  };

  // make keyword such as "#k1 #k2" to link to search query
  const parseKeyword = (keyword: string) => {
    // if keyword does not have space, return it directly
    if (keyword === undefined) {
      return null;
    }
    if (keyword.indexOf(' ') === -1) {
      return null;
    }
    const keywords = keyword.split(' ');
    // do not split them to multiple div element
    return (
      keywords.map((keyword, index) => (
        <Tag
          key={index.toString()}
          onClick={(event) => {
            makeNewSearch(keyword);
            event.stopPropagation();
          }}
          style={{
            cursor: 'pointer',
            marginRight: '4px',
            marginBottom: '4px',
            backgroundColor: 'rgba(var(--gray-6), 0.4)',
          }}
        >
          {keyword}
        </Tag>
      ))
    );
  };

  const className = cs(styles['card-block'], styles[`${type}-card`], styles[`zoom`]);


  return type === 'video' ? (
    <Card
      bordered={true}
      className={className}
      size="small"
      onClick={() => {goToVideoPage(card.video_uid)}}
      // cover_url as background image, width 100% and height 100%
      style={{
        backgroundImage: `url(${card.cover_url})`,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat',
        borderRadius: '8px',
      }}
    >
      <div className={styles['card-extra-like']}>
        <IconButton
          icon={<Like theme="filled" size="24" fill={card.is_user_liked?"red":"#ffffff"} onClick={()=> {console.log(card)}}/>}
        />
        { /* if liked, show red count text */ }
        <div className={styles['card-extra-like-count']} style={{ color: card.is_user_liked?'red':'#ffffff' }}>{card.like_count}</div>
      </div>
      {/* if saw before, show tag */}
      {card.is_user_watched ? (
        <div className={styles['card-extra-seen']}>
          <Tag
            icon={<IconClockCircle />}
            style={{
              backgroundColor: 'rgba(var(--gray-8), 0.5)',
            }}
          >观看过</Tag>
        </div>
      ) : null}
      {/* todo: if uploaded by this user, show video control */}

      <div className={styles['card-block-mask']}>
        {/*<div style={{ marginTop: '280px' }}></div>*/}
        <div
          className={cs(styles.title, {
            [styles['title-more']]: visible,
          })}
        >
          <div style={{ display: 'flex' }} onClick={(event) => {
            router.push({
              pathname: `/user/${card.user.username}`,
            });
            event.stopPropagation();
          }}>
            { /* add avatar to the left */}
            <Avatar size={40} style={{ marginTop: '4px' }}>
              {card.user?(card.user.avatar_url?<img src={card.user.avatar_url} />:card.user.nickname):'A'}
            </Avatar>
            <div style={{
              marginLeft: '8px',
            }}>
              <div className={styles.nickname}>{card.user?card.user.nickname:''}</div>
              <div className={styles.username}>@{card.user?card.user.username:''}</div>
            </div>
          </div>

          <div className={styles.content} >{card.content}</div>
          <div className={styles.keyword}>{parseKeyword(card.keyword)}</div>
          <div className={styles.time}>{parseTime(card.upload_time)}</div>
        </div>
      </div>
    </Card>
  ) : (
    <Card
      bordered={true}
      className={className}
      size="small"
    >
      <div
        className={cs(styles.title, {
          [styles['title-more']]: visible,
        })}
      >
        <div
          className={styles.usercard}
          style={{ display: 'flex' }}
          onClick={(event) => {
            router.push({
              pathname: `/user/${card.username}`,
            });
            event.stopPropagation();
          }}
        >
          { /* add avatar to the left */}
          <Avatar size={50} style={{ marginTop: '0' }}>
            {card?(card.avatar_url?<img src={card.avatar_url} />:card.nickname):'A'}
          </Avatar>
          <div style={{ marginLeft: '8px' }}>
            <div className={styles.nickname}>{card?card.nickname:''}</div>
            <div className={styles.username}>@{card?card.username:''}</div>
          </div>
        </div>
      </div>
    </Card>
  );
}

export default CardBlock;
