import React from 'react';
import styles from './style/index.module.less';

interface UserAddonCountInfoType {
  type: string;
  data: number;
}

export default function UserAddonCountInfo(props: UserAddonCountInfoType) {
  const { type, data} = props;

  function parseData(data: number) {
    if (data < 10000) {
      return data;
    }
    return (data / 10000).toFixed(1) + '万';
  }

  return (
    <div className={styles['user-addon-count-info']}>
      <div className={styles['user-addon-count-type']}>{type}</div>
      <div className={styles['user-addon-count-data']}>{parseData(data)}</div>
    </div>
  )
}