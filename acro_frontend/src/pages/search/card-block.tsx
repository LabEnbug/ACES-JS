import React, { useEffect, useRef, useState } from 'react';
import cs from 'classnames';
import { Button, Tag, Card, Avatar, Divider } from '@arco-design/web-react';
import styles from './style/index.module.less';
import { useRouter } from 'next/router';
import { Like } from '@icon-park/react';
import {
  IconCheck,
  IconEye,
  IconHeartFill,
  IconMinusCircle,
  IconPlus,
} from '@arco-design/web-react/icon';
import baxios from "@/utils/getaxios";
import useLocale from '@/utils/useLocale';
import locale from './locale';
import {parseTime} from "@/utils/timeUtils";
import {parseKeyword} from "@/utils/keywordUtils";

function CardBlock(props) {
  const t = useLocale(locale);
  const tg = useLocale();
  const { type, card = {} } = props;
  const [followLoading, setFollowLoading] = useState(false);
  const [followHovering, setFollowHovering] = useState(false);
  const [isFollowed, setIsFollowed] = useState(false);

  const router = useRouter();

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


  const followUser = (follow) => {
    setFollowLoading(true);
    // sleep 1000ms
    setTimeout(() => {
      (follow ? baxios.delete : baxios.post)
      ('/users/' + card.username + '/follow')
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
                  {t['cardBlock.video.liked']}
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
                  {t['cardBlock.video.watched']}
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
                {t['cardBlock.video.followed']}
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
                    <img src={card.user.avatar_url} alt={null} />
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
                {isFollowed ? (
                  (followHovering ? t['cardBlock.user.cancel'] : t['cardBlock.user.already']) + t['cardBlock.user.followed']
                ) : (t['cardBlock.user.follow'])}
              </Button>
            )}
          </div>
          <div className={styles['user-addon-info']}>
            <div className={styles['user-addon-count-info']}>
              <div className={styles['user-addon-count-type']}>{t['cardBlock.user.follower']}</div>
              <div className={styles['user-addon-count-data']}>
                {parseData(card ? card.be_followed_count : 0)}
              </div>
            </div>
            <Divider
              className={styles['user-addon-info-divider']}
              type="vertical"
            />
            <div className={styles['user-addon-count-info']}>
              <div className={styles['user-addon-count-type']}>{t['cardBlock.user.liked']}</div>
              <div className={styles['user-addon-count-data']}>
                {parseData(card ? card.be_liked_count : 0)}
              </div>
            </div>
            <Divider
              className={styles['user-addon-info-divider']}
              type="vertical"
            />
            <div className={styles['user-addon-count-info']}>
              <div className={styles['user-addon-count-type']}>{t['cardBlock.user.view']}</div>
              <div className={styles['user-addon-count-data']}>
                {parseData(card ? card.be_watched_count : 0)}
              </div>
            </div>
          </div>
        </div>
      </div>
    </Card>
  );
}

export default CardBlock;
