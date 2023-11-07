import React, { useEffect, useRef, useState } from 'react';
import {
  Tabs,
  Card,
  Button,
  List,
} from '@arco-design/web-react';
import useLocale from '@/utils/useLocale';
import locale from './locale';
import styles from './style/index.module.less';
import CardBlock from './card-block';
import { useRouter } from 'next/router';
import { Empty } from '@arco-design/web-react';
import { IconLoading } from '@arco-design/web-react/icon';
import baxios from "@/utils/getaxios";
import Head from "next/head";


const defaultVideoList = new Array(0).fill({});
const defaultUserList = new Array(0).fill({});
export default function ListSearchResult() {
  const t = useLocale(locale);
  const tg = useLocale();
  const [loading, setLoading] = useState(true);
  const [videoData, setVideoData] = useState(defaultVideoList);
  const [userData, setUserData] = useState(defaultUserList);

  const [activeKey, setActiveKey] = useState('video');

  const router = useRouter();
  const { q } = router.query;

  const listRef = useRef(null);

  const [isEndData, setIsEndData] = useState(false);

  const getData = (q, t) => {
    t === 'video'
      ? setVideoData(defaultVideoList)
      : setUserData(defaultUserList);
    setIsEndData(false);
    setLoading(true);
    // sleep
    // await new Promise(resolve => setTimeout(resolve, 3000));
    baxios
      .get(
        '/search/' + t + '?' +
        'keyword=' + encodeURIComponent(q) + '&' +
        'limit=' + '12'
      )
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          // 20231106 fixed
          setVideoData(defaultVideoList);
          setUserData(defaultUserList);
          return;
        }
        // 20231106 fixed: old result show when re-turn to the old tab
        if (t === 'video') {
          setVideoData(data.data.video_list);
          setUserData(defaultUserList);
        } else {
          setUserData(data.data.user_list);
          setVideoData(defaultVideoList);
        }
      })
      .catch((error) => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  };

  const getMoreData = (q, t) => {
    setLoading(true);
    baxios
      .get(
        '/search/' + t + '?' +
        'keyword=' + encodeURIComponent(q) + '&' +
        'start=' + (t === 'video' ? videoData.length : userData.length).toString() + '&' +
        'limit=' + '12'
      )
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          setIsEndData(true);
          return;
        }
        if (t === 'video') {
          setVideoData(videoData.concat(data.data.video_list));
        } else {
          setUserData(userData.concat(data.data.user_list));
        }
      })
      .catch((error) => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    if (router.isReady && q) {
      getData(q, activeKey);
    } else {
      window.location.href = '/';
    }
  }, [router.isReady, q, activeKey]);

  useEffect(() => {
    // add search history to local storage
    const searchHistory = localStorage.getItem('searchHistory');
    const maxHistoryNum = 10;
    const qTrim = (q as string).trim();
    if (searchHistory) {
      const historyList = JSON.parse(searchHistory);
      // todo: what if query q contains two # like "#k1 #k2"? or contains such as "#k1 value1", "value1 #k1", should we split it?

      if (historyList.indexOf(qTrim) === -1) {
        // put to the first
        historyList.unshift(qTrim);
      } else {
        // remove and put to the first
        const index = historyList.indexOf(qTrim);
        historyList.splice(index, 1);
        historyList.unshift(qTrim);
      }
      if (historyList.length > maxHistoryNum) {
        historyList.pop();
      }
      localStorage.setItem('searchHistory', JSON.stringify(historyList));
    } else {
      localStorage.setItem('searchHistory', JSON.stringify([qTrim]));
    }
  }, [router.isReady, q]);

  const ContentContainer = ({ children }) => (
    <div style={{ textAlign: 'center', marginTop: 4 }}>
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        {children}
      </div>
    </div>
  );

  const LoadingIndicator = () => (
    <div style={{ textAlign: 'center', marginTop: 36 }}>
      <span style={{ color: 'var(--color-text-3)' }}>
      </span>
    </div>
  );

  return (
    <>
      <Head>
        <title>{t['title']} - {tg['title.global']}</title>
      </Head>
      <Card className={styles['container']}>
        <Tabs
          activeTab={activeKey}
          type="rounded"
          onChange={setActiveKey}
          extra={
            loading && (
              <div style={{ color: 'var(--color-text-2)' }}>
                <IconLoading
                  style={{ marginRight: 8, color: 'rgb(var(--arcoblue-6))' }}
                />
                {t['search.loading']}
              </div>
            )
          }
        >
          <Tabs.TabPane key="video" title={t['cardList.tab.title.video']} />
          <Tabs.TabPane key="user" title={t['cardList.tab.title.user']} />
        </Tabs>
        <List
          ref={listRef}
          className={styles['card-list']}
          grid={
            activeKey === 'video'
              ? {
                  xs: 12,
                  sm: 12,
                  md: 12,
                  lg: 8,
                  xl: 6,
                  xxl: 4,
                }
              : {
                  xs: 12,
                  sm: 12,
                  md: 12,
                  lg: 12,
                  xl: 8,
                  xxl: 6,
                }
          }
          noDataElement={
            loading ? (
              <div />
            ) : (
              <Empty
                description={
                  <ContentContainer>
                    <span
                      style={{
                        color: 'var(--color-text-3)',
                        marginTop: '16px',
                        marginBottom: '16px',
                      }}
                    >
                      {t['search.noResultAbout']}{activeKey === 'video' ? t['search.video'] : t['search.user']}
                    </span>
                    <Button
                      type="text"
                      onClick={() => {
                        getData(q, activeKey);
                      }}
                    >
                      {t['search.searchAgain.before']}&quot;{q}&quot;{t['search.searchAgain.after']}
                    </Button>
                  </ContentContainer>
                }
              ></Empty>
            )
          }
          style={{ overflowY: 'scroll', height: 'calc(100vh - 170px)' }}
          dataSource={activeKey === 'video' ? videoData : userData}
          bordered={false}
          onReachBottom={() => {
            getMoreData(q, activeKey);
          }}
          render={(item, _) => (
            <List.Item style={{ padding: '4px 4px' }}>
              <CardBlock card={item} type={activeKey} loading={loading} />
            </List.Item>
          )}
          loading={loading}
          offsetBottom={300}
          footer={
            loading ? (
              <LoadingIndicator />
            ) : isEndData ? (
              <ContentContainer>
                <span
                  style={{ color: 'var(--color-text-3)', marginBottom: '4px' }}
                >
                  {t['search.noMoreContent']}
                </span>
                <Button
                  type="text"
                  onClick={() => {
                    getData(q, activeKey);
                  }}
                >
                  {t['search.searchAgain.before']}&quot;{q}&quot;{t['search.searchAgain.after']}
                </Button>
              </ContentContainer>
            ) : (
              (activeKey === 'video' ? videoData : userData).length !== 0 && (
                <ContentContainer>
                  <Button
                    type="text"
                    onClick={() => {
                      getMoreData(q, activeKey);
                    }}
                  >
                    {t['search.getMore']}{activeKey === 'video' ? t['search.video'] : t['search.user']}
                  </Button>
                </ContentContainer>
              )
            )
          }
        />
      </Card>
    </>
  );
}
