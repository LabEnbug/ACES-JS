import React, { forwardRef } from 'react';
import { Tag } from '@arco-design/web-react';
import styles from './style/intro.module.less';
import locale from './locale';
import useLocale from '@/utils/useLocale';


const parseKeyword = (keyword: string) => {
    // if keyword does not have space, return it directly
    if (keyword === undefined) {
        return null;
    }
    if (keyword.indexOf(' ') === -1) {

    }
    const keywords = keyword.split(' ');
    // do not split them to multiple div element
    return (
        keywords.map((keyword, index) => (
        <Tag
            key={index.toString()}
            style={{
            cursor: 'pointer',
            marginRight: '4px',
            marginBottom: '4px'}}
        >
            {keyword}
        </Tag>
        ))
    );
};

function BriefIntro(props, ref) {
  const { videoinfo, ...rest } = props;
  const t = useLocale(locale);
  return (
    <div className={styles['intro-contaner']}>
        <div className={styles['title-container ']}>
            <span className={styles['title']}>@{videoinfo['nickname']}</span>・<span className={styles['title-time']}>{videoinfo['time'].split('T')[0].replace(/-/g, '/')}</span>
        </div>
        <div className={styles['brief-container']}><span className={styles['brief-text']}> {videoinfo['content']} </span></div>
        {videoinfo['keyword'].includes('#') ? <div className={styles.keyword}>{parseKeyword(videoinfo['keyword'])}</div> : <></>}
    </div>
 );
}

export default forwardRef(BriefIntro);
