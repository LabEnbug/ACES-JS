import React, { forwardRef } from 'react';
import { Button, Tooltip, Space} from '@arco-design/web-react';
import styles from './style/index.module.less';
import cs from 'classnames';
import { Like, MessageUnread, Star, ShareTwo, More} from '@icon-park/react'
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

function LikeButton(props, ref) {
  const { icon, className, ...rest } = props;
  const t = useLocale(locale);
  return (
    <Space className={styles['icon-group']}
           direction='vertical'
           size={40}
    >
      <IconButton 
        icon={
          <>
          <Like theme="filled" size="36" fill="#ffffff" onClick={()=> {console.log('asdad')}}/>
          <p >4728</p>
          </>
        }
        tooltip={t['tooltip.like']}
      />
      <IconButton 
        icon={
          <>
            <MessageUnread theme="filled" size="36" fill="#ffffff"/>
            <p >4728</p>
          </>
        }
        tooltip={t['tooltip.comment']}
      />
      <IconButton 
        icon={
          <>
            <Star theme="outline" size="36" fill="#ffffff"/>
            <p >4728</p>
          </>
        }
        tooltip={t['tooltip.collection']}
      />
      <IconButton 
        icon={
          <>
            <ShareTwo theme="filled" size="36" fill="#ffffff"/>
            <p >4728</p>
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

export default forwardRef(LikeButton);
