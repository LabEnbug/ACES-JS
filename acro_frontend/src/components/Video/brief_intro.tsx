import React, { forwardRef } from 'react';
import { Tag } from '@arco-design/web-react';
import styles from './style/intro.module.less';
import locale from './locale';
import useLocale from '@/utils/useLocale';
import { useRouter } from 'next/router';
import {parseKeywordOnVideo} from "@/utils/keywordUtils";
import {parseTime} from "@/utils/timeUtils";

function BriefIntro(props, ref) {
  const { videoinfo, ...rest } = props;
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
        ・
        <span className={styles['title-time']}>{parseTime(videoinfo['time'], tg)}</span>
      </div>
      <div className={styles['brief-container']}><span className={styles['brief-text']}> {videoinfo['content']} </span></div>
      <div className={styles.keyword}>{parseKeywordOnVideo(videoinfo['keyword'], router)}</div>
    </div>
 );
}

export default forwardRef(BriefIntro);
