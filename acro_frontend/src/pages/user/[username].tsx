import React, {useEffect, useRef, useState} from 'react';
import axios from 'axios';
import {Tabs, Card, Input, Typography, Grid, Button, List, Divider, Avatar} from '@arco-design/web-react';
import useLocale from '@/utils/useLocale';
import locale from './locale';
import styles from './style/index.module.less';
import CardBlock from './card-block';
import {VideoCard} from './interface';
import {useRouter} from "next/router";
import {GlobalState} from "@/store";
import {Empty} from "@arco-design/web-react";
import {Popconfirm} from "@arco-design/web-react";
import {IconLoading} from "@arco-design/web-react/icon";
import active from "@antv/g2/src/interaction/action/element/active";
import GetAxios from "@/utils/getaxios";

const { Title } = Typography;
const { Row, Col } = Grid;

const defaultVideoList = new Array(0).fill({});
export default function ListSearchResult() {
  const t = useLocale(locale);
  const [loading, setLoading] = useState(true);
  const [videoData, setVideoData] = useState(defaultVideoList);
  const [videoNum, setVideoNum] = useState({});

  const [userData, setUserData] = useState(null);

  const [nicknameForChange, setNicknameForChange] = useState('');

  const [activeKey, setActiveKey] = useState('uploaded');

  const router = useRouter();
  const { username } = router.query;

  const listRef = useRef(null);

  const [isEndData, setIsEndData] = useState(false);

  const isSelf = username === 'self';

  const [noSuchUser, setNoSuchUser] = useState(false);

  const getUserInfoData = async (username) => {
    setLoading(true);
    const baxios = GetAxios();
    let params = new FormData();
    params.append('username', username);
    baxios.post('/v1-api/v1/user/query', params)
      .then(response => {
        const data = response.data
        if (data.status !== 200) {
          console.error(data.err_msg);
          setNoSuchUser(true);
          return;
        }
        setUserData(data.data.user);
        getVideoData(data.data.user.user_id, activeKey);
      })
      .catch(error => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  };

  const getSelfInfoData = async () => {
    setLoading(true);
    const baxios = GetAxios();
    baxios.post('/v1-api/v1/user/info')
      .then(response => {
        const data = response.data
        if (data.status !== 200) {
          console.error(data.err_msg);
          window.location.href = '/'; // reject not logged in
          return;
        }
        setUserData(data.data.user);
        getVideoData(data.data.user.user_id, activeKey);
      })
      .catch(error => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  }

  const getVideoData = async (userid, t) => {
    setIsEndData(false);
    setLoading(true);
    let param = new FormData();
    param.append('user_id', userid);
    param.append('action_history', t);
    // sleep
    // await new Promise(resolve => setTimeout(resolve, 3000));
    const baxios = GetAxios();
    baxios.post('/v1-api/v1/video/list' , param)
      .then(response => {
        const data = response.data
        if (data.status !== 200) {
          console.error(data.err_msg);
          setVideoData(defaultVideoList);
          return;
        }
        setVideoData(data.data.video_list);
        if (data.data.video_num) {
          let tmp = videoNum;
          tmp[t] = data.data.video_num;
          setVideoNum(tmp);
        }
      })
      .catch(error => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  };

  const getMoreData = async (userid, t) => {
    setLoading(true);
    let param = new FormData();
    param.append('user_id', userid);
    param.append('action_history', t);
    let s = videoData.length;
    param.append('start', s.toString());
    const baxios = GetAxios();
    baxios.post('/v1-api/v1/video/list', param)
      .then(response => {
        const data = response.data
        // sleep 1000ms
        setTimeout(() => {}, 3000);
        if (data.status !== 200) {
          console.error(data.err_msg);
          setIsEndData(true);
          return;
        }
        setVideoData(videoData.concat(data.data.video_list));
        if (data.data.video_num) {
          let tmp = videoNum;
          tmp[t] = data.data.video_num;
          setVideoNum(tmp);
        }
      })
      .catch(error => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  }

  useEffect(() => {
    if (router.isReady && username) {
      setVideoData(defaultVideoList);
      setUserData(null);
      (isSelf ? getSelfInfoData() : getUserInfoData(username));
    }
  }, [router.isReady, username]);

  useEffect(() => {
    // get video after user data set
    if (userData) {
      setVideoData(defaultVideoList);
      getVideoData(userData.user_id, activeKey);
    }
  }, [activeKey]);



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

  const ContentContainer = ({ children }) => (
    <div style={{ textAlign: 'center', marginTop: 4 }}>
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        {children}
      </div>
    </div>
  );

  const LoadingIndicator = () => (
    <div style={{ textAlign: 'center', marginTop: 16 }}>
    <span style={{ color: 'var(--color-text-3)' }}>
      <IconLoading style={{ marginRight: 8, color: 'rgb(var(--arcoblue-6))' }} />
        加载中
    </span>
    </div>
  );

  function handleChangeNickname() {
    return new Promise((resolve, reject) => {
      const baxios = GetAxios();
      let params = new FormData();
      params.append('nickname', nicknameForChange);
      baxios.post('/v1-api/v1/user/info/set', params)
        .then(response => {
          const data = response.data
          if (data.status !== 200) {
            console.error(data.err_msg);
            return;
          }
          setUserData(data.data.user);
          resolve();
        })
        .catch(error => {
          console.error(error);
          reject();
        })
        .finally(() => {});
    });
  }

  return (
    <>
      <Card className={styles['info-wrapper']}>
        { noSuchUser ? (
          <div style={{ textAlign: 'center', marginTop: 16 }}>
            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
              <span style={{ color: 'var(--color-text-1)' }}>
                用户不存在
              </span>
            </div>
          </div>
        ) : (
          <div style={{ display: 'flex' }}>
            { /* add avatar to the left */}
            <Avatar size={64} style={{ }}>
              {userData?(userData.avatar_url?<img src={userData.avatar_url} />:userData.nickname):'A'}
            </Avatar>
            <div style={{
              marginLeft: '16px',
              //vertical center
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'center',
            }}>
              <div className={styles.nickname}>
                {userData?userData.nickname:''}
                {isSelf?(
                    <Popconfirm
                      position='bottom'
                      icon={null}
                      title={<Input placeholder='请输入新的昵称' value={nicknameForChange} onChange={(e) => setNicknameForChange(e)}/>}
                      okText='修改'
                      cancelText='取消'
                      onOk={handleChangeNickname}
                      onCancel={() => {console.log('cancel')}}
                    >
                      <Button type="text" className={styles['change-nickname']}>修改昵称</Button>
                    </Popconfirm>
                ):null}
              </div>
              <div className={styles.username}>@{userData?userData.username:''}</div>
            </div>
          </div>
        )}
      </Card>
      { noSuchUser ? null : (
        <>
          <div style={{ marginTop: '16px' }}></div>
          <Card
            // style={{ height: 'calc(100vh - 150px)' }}
          >
            {/*<Title heading={6}>{t['menu.list.card']}</Title>*/}
            <Tabs
              activeTab={activeKey}
              type="text"
              onChange={setActiveKey}
            >
              <Tabs.TabPane key="uploaded" title={t['cardList.tab.title.uploaded'] + (videoNum['uploaded']&&videoNum['uploaded']!==0?" ("+videoNum['uploaded']+")":"")} />
              <Tabs.TabPane key="like" title={t['cardList.tab.title.like'] + (videoNum['like']&&videoNum['like']!==0?" ("+videoNum['like']+")":"")} />
              <Tabs.TabPane key="favorite" title={t['cardList.tab.title.favorite'] + (videoNum['favorite']&&videoNum['favorite']!==0?" ("+videoNum['favorite']+")":"")} />
              {/* history is only for self */ }
              {isSelf ? <Tabs.TabPane key="history" title={t['cardList.tab.title.history'] + (videoNum['history']&&videoNum['history']!==0?" ("+videoNum['history']+")":"")} /> : null}
            </Tabs>
            <Divider />
            <List
              ref={listRef}
              grid={{
                sm: 24,
                md: 12,
                lg: 8,
                xl: 8,
              }}
              noDataElement={loading?<div />:<Empty
                description={ ' ' }
              ></Empty>}
              dataSource={videoData}
              bordered={false}
              onListScroll={() => {console.log(1111)}}
              onReachBottom={() => {console.log(111)}}
              render={(item, index) => (
                <List.Item style={{ padding: '4px 4px' }}>
                  <CardBlock card={item} type={activeKey} watching_username={userData?userData.username:''} loading={loading}/>
                </List.Item>
              )}
            />
            { loading ? (
              <LoadingIndicator />
            ) : (
              isEndData ? (
                <ContentContainer>
                  <span style={{ color: 'var(--color-text-3)', marginBottom: '4px' }}>无更多内容</span>
                </ContentContainer>
              ) : (
                videoData.length === 0 ? (
                  <ContentContainer>
            <span style={{color: 'var(--color-text-3)', marginBottom: '4px',}}>
              { /* uploaded, like, favorite, history info */ }
              {`${isSelf ? '您' : '该用户'}还没有${activeKey === 'uploaded' ? '上传' : activeKey === 'like' ? '点赞' : activeKey === 'favorite' ? '收藏' : '观看'}过任何视频`}
            </span>
                  </ContentContainer>
                ) : (
                  <ContentContainer>
                    <Button type='text' onClick={() => { getMoreData(userData.user_id, activeKey) }}>
                      {`加载更多${activeKey === 'uploaded' ? '上传' : activeKey === 'like' ? '点赞' : activeKey === 'favorite' ? '收藏' : '观看'}的视频`}
                    </Button>
                  </ContentContainer>
                )
              )
            )
            }
          </Card>
        </>
      )
      }
    </>
  );
}
