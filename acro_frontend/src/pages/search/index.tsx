import React, {useEffect, useRef, useState} from 'react';
import axios from 'axios';
import {Tabs, Card, Input, Typography, Grid, Button, List} from '@arco-design/web-react';
import useLocale from '@/utils/useLocale';
import locale from './locale';
import styles from './style/index.module.less';
import CardBlock from './card-block';
import {VideoCard} from './interface';
import {useRouter} from "next/router";
import {GlobalState} from "@/store";
import {Empty} from "antd";
import {IconLoading} from "@arco-design/web-react/icon";
import active from "@antv/g2/src/interaction/action/element/active";

const { Title } = Typography;
const { Row, Col } = Grid;

const defaultVideoList = new Array(0).fill({});
const defaultUserList = new Array(0).fill({});
export default function ListSearchResult() {
  const t = useLocale(locale);
  const [loading, setLoading] = useState(true);
  const [videoData, setVideoData] = useState(defaultVideoList);
  const [userData, setUserData] = useState(defaultUserList);

  const [activeKey, setActiveKey] = useState('video');

  const router = useRouter();
  const { q } = router.query;

  const listRef = useRef(null);

  const getData = async (q, t) => {
    console.log(q)
    setLoading(true);
    let param = new FormData();
    param.append('keyword', q);
    // sleep
    // await new Promise(resolve => setTimeout(resolve, 3000));
    axios.post(t === 'video' ? '/v1-api/v1/video/search' : '/v1-api/v1/user/search', param)
      .then(response => {
        const data = response.data
        // sleep 1000ms
        setTimeout(() => {}, 3000);
        if (data.status !== 200) {
          console.error(data.err_msg);
          t === 'video' ? setVideoData([]) : setUserData([])
          return;
        }
        t === 'video' ? setVideoData(data.data.video_list) : setUserData(data.data.user_list)
      })
      .catch(error => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  };

  const getMoreData = async (q, t, s) => {
    console.log(q, t, s)
    setLoading(true);
    let param = new FormData();
    param.append('keyword', q);
    param.append('start', s);
    axios.post(t === 'video' ? '/v1-api/v1/video/search' : '/v1-api/v1/user/search', param)
      .then(response => {
        const data = response.data
        // sleep 1000ms
        setTimeout(() => {}, 3000);
        if (data.status !== 200) {
          console.error(data.err_msg);
          return;
        }
        if (t === 'video') {
          setVideoData(videoData.concat(data.data.video_list))
        } else {
          setUserData(userData.concat(data.data.user_list))
        }
      })
      .catch(error => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  }

  useEffect(() => {
    if (router.isReady && q) {
      activeKey === 'video' ? setVideoData(defaultVideoList) : setUserData(defaultUserList);
      getData(q, activeKey);
      // add search history to local storage
      const searchHistory = localStorage.getItem('searchHistory');
      const maxHistoryNum = 10;
      if (searchHistory) {
        const historyList = JSON.parse(searchHistory);
        if (historyList.indexOf(q) === -1) {
          // put to the first
          historyList.unshift(q);
        } else {
          // remove and put to the first
          const index = historyList.indexOf(q);
          historyList.splice(index, 1);
          historyList.unshift(q);
        }
        if (historyList.length > maxHistoryNum) {
          historyList.pop();
        }
        localStorage.setItem('searchHistory', JSON.stringify(historyList));
      } else {
        localStorage.setItem('searchHistory', JSON.stringify([q]));
      }
    }
  }, [router.isReady, q, activeKey]);

  useEffect(() => {
    if (!loading && listRef.current) {
      const listElement = listRef.current;
      // 检查内容高度是否小于等于容器高度
      console.log(11)
      if (listElement.clientHeight >= listElement.scrollHeight) {
        // 触发到达底部的逻辑
        console.log('Reached bottom');
        // 在这里调用onReachBottom或相关逻辑
      }
    }
  }, [loading]);


  const getCardList = (
    list: Array<VideoCard>,
    type: 'video' | 'user'
  ) => {
    if (list.length === 0) {
      return (
        <Empty
          className={styles['empty']}
          image={Empty.PRESENTED_IMAGE_DEFAULT}
          description={ '搜索不到相关的' + (type === 'video' ? '视频' : '用户') + '，请尝试其他关键词' }
        ><Button type='text' onClick={() => { getData(q, activeKey) }}>重新尝试搜索</Button></Empty>
      );
    } else {
      return (
        <Row gutter={24} className={styles['card-content']}>
          {list.map((item, index) => (
            <Col xs={24} sm={12} md={8} lg={8} xl={8} xxl={8} key={index}>
              <CardBlock card={item} type={type} loading={loading}/>
            </Col>
          ))}
        </Row>
      );
    }
  };

  return (
    <Card
      // style={{ height: 'calc(100vh - 150px)' }}
    >
      {/*<Title heading={6}>{t['menu.list.card']}</Title>*/}
      <Tabs
        activeTab={activeKey}
        type="rounded"
        onChange={setActiveKey}
      >
        <Tabs.TabPane key="video" title={t['cardList.tab.title.video']} />
        <Tabs.TabPane key="user" title={t['cardList.tab.title.user']} />
      </Tabs>
      <List
        ref={listRef}
        grid={{
          sm: 24,
          md: 12,
          lg: 8,
          xl: 8,
        }}
        noDataElement={loading?<div />:<Empty
          className={styles['empty']}
          image={Empty.PRESENTED_IMAGE_DEFAULT}
          description={ '搜索不到相关的' + (activeKey === 'video' ? '视频' : '用户') + '，请尝试其他关键词' }
        ><Button type='text' onClick={() => { getData(q, activeKey) }}>重新尝试搜索</Button></Empty>}
        dataSource={activeKey === 'video' ? videoData : userData}
        bordered={false}
        onListScroll={() => {console.log(1111)}}
        onReachBottom={() => {console.log(111)}}
        render={(item, index) => (
          <List.Item style={{ padding: '4px 4px' }}>
            <CardBlock card={item} type={activeKey} loading={loading}/>
          </List.Item>
        )}
      />
      { loading &&
      <div style={{ textAlign: 'center', marginTop: 16 }}>
        <span style={{ color: 'var(--color-text-3)' }}>
          <IconLoading style={{ marginRight: 8, color: 'rgb(var(--arcoblue-6))' }} />
            加载中
        </span>
      </div>
      }
    </Card>
  );
}
