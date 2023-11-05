import React, { useEffect, useState } from 'react';
import cs from 'classnames';
import {
  Tag,
  Card,
  Dropdown,
  Menu,
  Avatar,
  Popconfirm,
  Message,
} from '@arco-design/web-react';
import useLocale from '@/utils/useLocale';
import locale from './locale';
import { VideoCard } from './interface';
import styles from './style/index.module.less';
import { useRouter } from 'next/router';
import { Like } from '@icon-park/react';
import {
  IconDelete,
  IconEdit,
  IconEye, IconFire,
  IconHeartFill,
  IconLiveBroadcast, IconLock,
  IconMore,
  IconToTop, IconUnlock,
} from '@arco-design/web-react/icon';
import baxios from "@/utils/getaxios";

interface CardBlockType {
  type: string;
  card: VideoCard;
  watching_username: string;
  loading?: boolean;
  onDelete: () => void;
}

function CardBlock(props: CardBlockType) {
  const { type, card, watching_username, onDelete } = props;
  const [isVideoCardPopup, setIsVideoCardPopup] = useState(false);
  const [deleteLoading, setDeleteLoading] = useState(false);

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

  function topVideo() {
    const params = new FormData();
    params.append('video_uid', card.video_uid);
    params.append('type', card.is_top ? 'untop' : 'top');
    baxios
      .post('/v1-api/v1/video/top', params)
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          Message.error({
            content: '短视频' + (card.is_top ? '取消' : '') + '置顶失败！',
            duration: 5000,
          });
          return;
        }
        Message.success({
          content: '短视频' + (card.is_top ? '取消' : '') + '置顶成功！',
          duration: 5000,
        });
        card.is_top = !card.is_top;
      })
      .catch((error) => {
        console.error(error);
      })
      .finally(() => {
        if (isVideoCardPopup) {
          setIsVideoCardPopup(false);
        }
      });
  }

  function privateVideo() {
    const params = new FormData();
    params.append('video_uid', card.video_uid);
    params.append('type', card.is_private ? 'unprivate' : 'private');
    baxios
      .post('/v1-api/v1/video/private', params)
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          Message.error({
            content: '短视频' + (card.is_private ? '取消' : '') + '置顶失败！',
            duration: 5000,
          });
          return;
        }
        Message.success({
          content: '短视频' + (card.is_private ? '取消' : '') + '置顶成功！',
          duration: 5000,
        });
        card.is_private = !card.is_private;
      })
      .catch((error) => {
        console.error(error);
      })
      .finally(() => {
        if (isVideoCardPopup) {
          setIsVideoCardPopup(false);
        }
      });
  }

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

  return (
    <Card
      bordered={true}
      className={cs(styles['card-block'], styles[`video-card`], styles[`zoom`])}
      size="small"
      onClick={() => {
        if (!isVideoCardPopup) {
          goToVideoPage(card.video_uid);
        }
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
        {(card.is_private || card.is_top) && (
          <div style={{ display: 'flex', marginBottom: '8px' }}>
            {card.is_top && (
            <div className={styles['video-card-extra-seen']}>
              <Tag
                icon={<IconToTop />}
                style={{
                  backgroundColor: 'rgba(var(--gray-8), 0.5)',
                }}
              >
                置顶
              </Tag>
            </div>
            )}
            {card.is_private && (
              <div className={styles['video-card-extra-seen']}>
                <Tag
                  icon={<IconLock />}
                  style={{
                    backgroundColor: 'rgba(var(--gray-8), 0.5)',
                  }}
                >
                  私密
                </Tag>
              </div>
            )}
          </div>
        )}
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
            {card.is_user_watched && type !== 'watched' && (
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
      </div>
      <div className={styles['video-card-uploaded-control']}>
        {card.is_user_uploaded && (
          <Dropdown
            popupVisible={isVideoCardPopup}
            onVisibleChange={setIsVideoCardPopup}
            droplist={
              <Menu
                onClickMenuItem={(key) => {
                  return false;
                }}
              >
                <Menu.Item
                  key="makeTop"
                  onClick={(event) => {
                    event.stopPropagation();
                    topVideo();
                  }}
                >
                  <IconToTop className={styles['video-dropdown-icon']} />
                  视频{card.is_top&&'取消'}置顶
                </Menu.Item>
                <Menu.Item
                  key="edit"
                  onClick={(event) => {
                    event.stopPropagation();
                    router.push({
                      pathname: `/edit`,
                      query: {
                        video_uid: card.video_uid,
                      },
                    });
                  }}
                >
                  <IconEdit className={styles['video-dropdown-icon']} />
                  编辑视频信息
                </Menu.Item>
                <Menu.Item
                  key="promote"
                  onClick={(event) => {
                    event.stopPropagation();
                    router.push({
                      pathname: `/promote`,
                      query: {
                        video_uid: card.video_uid,
                      },
                    });
                  }}
                >
                  <IconFire
                    className={styles['video-dropdown-icon']}
                  />
                  推广流量购买
                </Menu.Item>
                <Menu.Item
                  key="advertise"
                  onClick={(event) => {
                    event.stopPropagation();
                    router.push({
                      pathname: `/advertise`,
                      query: {
                        video_uid: card.video_uid,
                      },
                    });
                  }}
                >
                  <IconLiveBroadcast
                    className={styles['video-dropdown-icon']}
                  />
                  投放广告
                </Menu.Item>
                <Menu.Item
                  key="makeTop"
                  onClick={(event) => {
                    event.stopPropagation();
                    privateVideo();
                  }}
                >
                  {card.is_private?(
                    <IconUnlock className={styles['video-dropdown-icon']} />
                  ):(
                    <IconLock className={styles['video-dropdown-icon']} />
                    )}{card.is_private?'取消':'设为'}私密
                </Menu.Item>
                <Menu.Item key="delete">
                  <Popconfirm
                    focusLock
                    title="确定要删除该视频吗？该操作不可撤销"
                    position="bottom"
                    cancelButtonProps={{
                      style: { marginRight: '8px' },
                      disabled: deleteLoading,
                    }}
                    okButtonProps={{ loading: deleteLoading }}
                    onOk={() => {
                      setDeleteLoading(true);
                      const params = new FormData();
                      params.append('video_uid', card.video_uid);
                      baxios
                        .post('/v1-api/v1/video/delete', params)
                        .then((response) => {
                          const data = response.data;
                          if (data.status !== 200) {
                            console.error(data.err_msg);
                            Message.error({
                              content: '短视频删除失败！',
                              duration: 5000,
                            });
                            return;
                          }
                          Message.success({
                            content: '短视频删除成功！',
                            duration: 5000,
                          });
                          onDelete();
                        })
                        .catch((error) => {
                          console.error(error);
                        })
                        .finally(() => {
                          setDeleteLoading(false);
                        });
                    }}
                  >
                    <IconDelete className={styles['video-dropdown-icon']} />
                    删除视频
                  </Popconfirm>
                </Menu.Item>
              </Menu>
            }
            position="bottom"
          >
            <Tag
              size={'large'}
              icon={<IconMore />}
              // onClick={(e) => e.stopPropagation()}
              style={{
                backgroundColor: 'rgba(var(--gray-8), 0.1)',
              }}
            ></Tag>
          </Dropdown>
        )}
      </div>
      {/* todo: if uploaded by this user, add video control button */}

      <div className={styles['video-card-bottom-mask']}>
        <div className={styles['video-card-bottom']}>
          <div
            className={styles['video-user-card-block']}
            onClick={(event) => {
              card.user.username !== watching_username &&
                router.push({
                  pathname: '/user/' + card.user.username,
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
  );
}

export default CardBlock;
