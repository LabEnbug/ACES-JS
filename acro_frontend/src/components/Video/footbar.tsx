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

function FootBar(props, ref) {
  const { playRef, visible } = props;
  const t = useLocale(locale);
  return (
    <></>
  );
}

export default forwardRef(FootBar);
