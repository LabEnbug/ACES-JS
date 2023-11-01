import React, { forwardRef } from 'react';
import { Tag } from '@arco-design/web-react';
import styles from './style/intro.module.less';
import locale from './locale';
import useLocale from '@/utils/useLocale';
import { useRouter } from 'next/router';

function BriefIntro(props, ref) {
  const { videoinfo, ...rest } = props;
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
    return (
      keywords.map((keyword, index) => (
        <Tag
          key={index.toString()}
          onClick={() => makeNewSearch(keyword)}
          style={{
            cursor: 'pointer',
            marginRight: '4px',
            marginBottom: '4px',
          }}
        >
          {keyword}
        </Tag>
      ))
    );
  };

  const goToUserPage = () => {
    router.push({
      pathname: `/user/${videoinfo['username']}`,
    });
  }
  return (
    <div className={styles['intro-contaner']}>
      <div className={styles['title-container']}>
        <>
          <span className={styles['title']} onClick={goToUserPage}>{videoinfo['nickname']}</span>
          <span className={styles['username']} onClick={goToUserPage}>@{videoinfo['username']}</span>
        </>
        ・<span className={styles['title-time']}>{videoinfo['time'].split('T')[0].replace(/-/g, '/')}</span>
      </div>
      <div className={styles['brief-container']}><span className={styles['brief-text']}> {videoinfo['content']} </span></div>
      <div className={styles.keyword}>{parseKeyword(videoinfo['keyword'])}</div>
    </div>
 );
}

export default forwardRef(BriefIntro);
