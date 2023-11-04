import React, { useEffect, useRef, useState } from 'react';
import cs from 'classnames';
import { Button, Tag, Card, Avatar, Divider } from '@arco-design/web-react';
import { VideoCard, UserCard } from './interface';
import styles from './style/index.module.less';
import { useRouter } from 'next/router';
import { Like } from '@icon-park/react';
import IconButton from '@/components/NavBar/IconButton';
import {
  IconCheck,
  IconEye,
  IconHeartFill,
  IconMinusCircle,
  IconPlus,
} from '@arco-design/web-react/icon';
import GetAxios from '@/utils/getaxios';

interface CardBlockType {
  type: 'video' | 'user';
  card?: VideoCard | UserCard;
  loading?: boolean;
}

function CardBlock(props) {
  const { type, card = {} } = props;
  const [visible, setVisible] = useState(false);

  const [loading, setLoading] = useState(props.loading);

  const [followLoading, setFollowLoading] = useState(false);
  const [followHovering, setFollowHovering] = useState(false);
  const [isFollowed, setIsFollowed] = useState(false);

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
    if (type === 'user') {
      setIsFollowed(card.be_followed);
    }
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
    return keywords.map((keyword, index) => (
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
    ));
  };

  const followUser = (follow) => {
    setFollowLoading(true);
    const baxios = GetAxios();
    const params = new FormData();
    params.append('user_id', card.user_id);
    params.append('action', follow ? 'unfollow' : 'follow');
    // sleep 1000ms
    setTimeout(() => {
      baxios
        .post('/v1-api/v1/user/follow', params)
        .then((response) => {
          const data = response.data;
          if (data.status !== 200) {
            console.error(data.err_msg);
            return;
          }
          setIsFollowed(!follow);
          card.be_followed = !follow;
        })
        .catch((error) => {
          console.error(error);
        })
        .finally(() => {
          setFollowLoading(false);
        });
    }, 1000);
  };

  function parseData(data: number) {
    if (data < 10000) {
      return data;
    }
    return (data / 10000).toFixed(1) + '万';
  }

  return type === 'video' ? (
    <Card
      bordered={true}
      className={cs(
        styles['card-block'],
        styles[`${type}-card`],
        styles[`zoom`]
      )}
      size="small"
      onClick={() => {
        goToVideoPage(card.video_uid);
      }}
      // cover_url as background image, width 100% and height 100%
      style={{
        backgroundImage: `url(${card.cover_url})`,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat',
        borderRadius: '8px',
      }}
    >
      <div className={styles['video-card-extra-like']}>
        <Like
          theme="filled"
          size="24"
          fill={card.is_user_liked ? 'red' : '#ffffff'}
          onClick={(event) => {
            console.log(card);
            event.stopPropagation();
          }}
        />
        <div
          className={styles['video-card-extra-like-count']}
          style={{ color: card.is_user_liked ? 'red' : '#ffffff' }}
        >
          {card.be_liked_count}
        </div>
      </div>

      {/* if saw before, show tag */}
      <div style={{ marginTop: '12px' }}>
        {(card.is_user_liked || card.is_user_watched) && (
          <div style={{ display: 'flex', marginBottom: '8px' }}>
            {card.is_user_liked && (
              <div className={styles['video-card-extra-seen']}>
                <Tag
                  icon={<IconHeartFill />}
                  style={{
                    backgroundColor: 'rgba(var(--gray-8), 0.5)',
                  }}
                >
                  点赞过
                </Tag>
              </div>
            )}
            {card.is_user_watched && (
              <div className={styles['video-card-extra-seen']}>
                <Tag
                  icon={<IconEye />}
                  style={{
                    backgroundColor: 'rgba(var(--gray-8), 0.5)',
                  }}
                >
                  观看过
                </Tag>
              </div>
            )}
          </div>
        )}

        {card.user.be_followed && (
          <div style={{ display: 'flex', marginBottom: '8px' }}>
            <div className={styles['video-card-extra-seen']}>
              <Tag
                style={{
                  backgroundColor: 'rgba(var(--gray-8), 0.5)',
                }}
              >
                关注的用户
              </Tag>
            </div>
          </div>
        )}
      </div>
      <div className={styles['video-card-bottom-mask']}>
        <div className={styles['video-card-bottom']}>
          <div
            className={styles['video-user-card-block']}
            onClick={(event) => {
              router.push({
                pathname: `/user/${card.user.username}`,
              });
              event.stopPropagation();
            }}
          >
            {/* add avatar to the left */}
            <div className={styles['user-card-left']}>
              <Avatar size={40}>
                {card.user ? (
                  card.user.avatar_url ? (
                    <img src={card.user.avatar_url} />
                  ) : (
                    card.user.nickname
                  )
                ) : (
                  'A'
                )}
              </Avatar>
            </div>
            <div className={styles['user-card-right']}>
              <div className={styles['user-name-info']}>
                <div className={styles.nickname}>
                  {card.user ? card.user.nickname : ''}
                </div>
                <div className={styles.username}>
                  @{card.user ? card.user.username : ''}
                </div>
              </div>
            </div>
          </div>
          <div className={styles.content}>{card.content}</div>
          <div className={styles.keyword}>{parseKeyword(card.keyword)}</div>
          <div className={styles.time}>{parseTime(card.upload_time)}</div>
        </div>
      </div>
    </Card>
  ) : (
    <Card
      bordered={true}
      className={cs(styles['card-block'], styles[`zoom`])}
      size="small"
    >
      <div
        className={styles['user-card-block']}
        onClick={(event) => {
          router.push({
            pathname: `/user/${card.username}`,
          });
          event.stopPropagation();
        }}
      >
        {/* add avatar to the left */}
        <div className={styles['user-card-left']}>
          <Avatar size={64}>
            {card.avatar_url ? <img src={card.avatar_url} /> : card.nickname}
          </Avatar>
        </div>
        <div className={styles['user-card-right']}>
          <div className={styles['user-card-right-top']}>
            <div className={styles['user-name-info']}>
              <div className={styles.nickname}>{card.nickname}</div>
              <div className={styles.username}>@{card.username}</div>
            </div>
            {!card.is_self && (
              <Button
                className={styles['user-follow']}
                size={'mini'}
                type={isFollowed ? 'secondary' : 'primary'}
                icon={
                  isFollowed ? (
                    followHovering ? (
                      <IconMinusCircle />
                    ) : (
                      <IconCheck />
                    )
                  ) : (
                    <IconPlus />
                  )
                }
                onClick={(e) => {
                  followUser(isFollowed);
                  e.stopPropagation();
                }}
                loading={followLoading}
                onMouseEnter={() => setFollowHovering(true)}
                onMouseLeave={() => setFollowHovering(false)}
              >
                {isFollowed && (followHovering ? '取消' : '已')}关注
              </Button>
            )}
          </div>
          <div className={styles['user-addon-info']}>
            <div className={styles['user-addon-count-info']}>
              <div className={styles['user-addon-count-type']}>粉丝</div>
              <div className={styles['user-addon-count-data']}>
                {parseData(card ? card.be_followed_count + 532422 : 0)}
              </div>
            </div>
            <Divider
              className={styles['user-addon-info-divider']}
              type="vertical"
            />
            <div className={styles['user-addon-count-info']}>
              <div className={styles['user-addon-count-type']}>获赞</div>
              <div className={styles['user-addon-count-data']}>
                {parseData(card ? card.be_liked_count + 11111 : 0)}
              </div>
            </div>
            <Divider
              className={styles['user-addon-info-divider']}
              type="vertical"
            />
            <div className={styles['user-addon-count-info']}>
              <div className={styles['user-addon-count-type']}>浏览量</div>
              <div className={styles['user-addon-count-data']}>
                {parseData(card ? card.be_watched_count + 12312 : 0)}
              </div>
            </div>
          </div>
        </div>
      </div>
    </Card>
  );
}

export default CardBlock;
