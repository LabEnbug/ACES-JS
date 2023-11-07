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
import {parseKeyword} from "@/utils/keywordUtils";
import {parseTime} from "@/utils/timeUtils";

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
  const tg = useLocale();

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
    (card.is_top ? baxios.delete : baxios.post)
    ('/videos/' + card.video_uid + '/actions/' + 'top')
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
    (card.is_private ? baxios.delete : baxios.post)
    ('/videos/' + card.video_uid + '/actions/' + 'private')
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          Message.error({
            content: '短视频' + (card.is_private ? '取消' : '') + '设为私密失败！',
            duration: 5000,
          });
          return;
        }
        Message.success({
          content: '短视频' + (card.is_private ? '取消' : '') + '设为私密成功！',
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
                  key="makePrivate"
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
                      baxios
                        .delete('/videos/' + card.video_uid)
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
          <div className={styles.keyword}>{parseKeyword(card.keyword, router)}</div>
          <div className={styles.time}>{parseTime(card.upload_time, tg)}</div>
        </div>
      </div>
    </Card>
  );
}

export default CardBlock;
