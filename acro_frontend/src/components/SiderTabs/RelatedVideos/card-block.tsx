import React, {useEffect} from 'react';
import cs from 'classnames';
import {
  Tag,
  Card,
} from '@arco-design/web-react';
import useLocale from '@/utils/useLocale';
import locale from './locale';
import { VideoCard } from './interface';
import styles from './style/index.module.less';
import { useRouter } from 'next/router';
import { Like } from '@icon-park/react';
import {
  IconEye,
  IconHeartFill,
} from '@arco-design/web-react/icon';

interface CardBlockType {
  card: VideoCard;
  loading?: boolean;
}

function CardBlock(props: CardBlockType) {
  const { card } = props;

  const t = useLocale(locale);
  const tg = useLocale();

  const router = useRouter();

  function goToVideoPage(video_uid: string) {
    router.push({
      pathname: `/video`,
      query: {
        type: "comprehensive",
        video_uid: video_uid,
        rel: 1,
      },
    }, undefined, {shallow:true});
    // window.location.reload();
    return;
  }

  return (
    <Card
      bordered={true}
      className={cs(styles['card-block'], styles[`zoom`])}
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
      <div style={{ marginTop: 4, marginLeft: 4 }}>
        {(card.is_user_watched) && (
          <div style={{ display: 'flex', marginBottom: '8px' }}>
            {card.is_user_watched && (
              <div className={styles['video-card-extra-seen']}>
                <Tag
                  icon={<IconEye />}
                  style={{
                    backgroundColor: 'rgba(var(--gray-8), 0.2)',
                  }}
                />
              </div>
            )}
          </div>
        )}
      </div>
      <div className={styles['video-card-bottom-like']}>
        <Like
          theme="filled"
          size="16"
          fill={card.is_user_liked ? 'red' : '#ffffff'}
          onClick={(event) => {
            console.log(card);
            event.stopPropagation();
          }}
        />
        <div
          className={styles['video-card-bottom-like-count']}
          style={{ color: card.is_user_liked ? 'red' : '#ffffff' }}
        >
          {card.be_liked_count}
        </div>
      </div>
    </Card>
  );
}

export default CardBlock;
