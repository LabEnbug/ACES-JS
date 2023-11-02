import React, { forwardRef } from 'react';
import { Button, Tooltip, Space, Avatar} from '@arco-design/web-react';
import styles from './style/index.module.less';
import cs from 'classnames';
import { Like, MessageUnread, Star, ShareTwo, More} from '@icon-park/react'
import { IconCheck, IconPlus } from '@arco-design/web-react/icon';
import locale from './locale';
import useLocale from '@/utils/useLocale';

function IconButton(props) {
  const { icon, tooltip } = props;

  return (
    <Tooltip position='lt' trigger='hover' content={tooltip}>
        <Button
        icon={icon}
        shape="square"
        type="secondary"
        className={cs(styles['icon-button'])}
    />
    </Tooltip>
  );
}

function SideBar(props, ref) {
  const { videoinfo, clicklike, className, clickfoward, userlike, ikecount, userfavorite, followed, changefollow, favoritecount, clickfavorite, forwardedcount, commentedcount, ...rest } = props;
  const t = useLocale(locale);
  return (
    <Space className={styles['icon-group']}
           direction='vertical'
           size={40}
    >
      <Tooltip position='lt' trigger='hover' content={ followed ? t['tooltip.unfollow'] : t['tooltip.follow'] }>
        <Avatar
          triggerIcon={ followed ? <IconCheck /> : <IconPlus />}
          triggerIconStyle={{
            color: followed ? 'rgba(255, 0, 0, 0.8)' : '#3491FA',
          }}
          onClick={changefollow}
          autoFixFontSize={true}
          style={{
            backgroundColor: '#000000',
          }}
        >
          {videoinfo.nickname}
        </Avatar>
      </Tooltip>
      <IconButton 
        icon={
          <>
          <Like theme="filled" size="36" fill={ userlike ? 'rgba(255, 0, 0, 0.8)' : "#ffffff"} onClick={()=>{clicklike['func'](...clicklike['params'])}}/>
          <p> {ikecount} </p>
          </>
        }
        tooltip={t['tooltip.like']}
      />
      <IconButton 
        icon={
          <>
            <MessageUnread theme="filled" size="36" fill="#ffffff"/>
            <p >{commentedcount}</p>
          </>
        }
        tooltip={t['tooltip.comment']}
      />
      <IconButton 
        icon={
          <>
            <Star theme="outline" size="36" fill={ userfavorite ? 'rgba(218, 165, 32, 0.8)' : "#ffffff"} onClick={()=>{clickfavorite['func'](...clickfavorite['params'])}} />
            <p >{favoritecount}</p>
          </>
        }
        tooltip={t['tooltip.favorite']}
      />
      <IconButton 
        icon={
          <>
            <ShareTwo theme="filled" size="36" fill="#ffffff" onClick={clickfoward}/>
            <p >{forwardedcount}</p>
          </>
        }
        tooltip={t['tooltip.forward']}
      />
      <IconButton icon={
        <>
          <More theme="filled" size="36" fill="#ffffff"/>
        </>
      }/>
    </Space>
  );
}

export default forwardRef(SideBar);
