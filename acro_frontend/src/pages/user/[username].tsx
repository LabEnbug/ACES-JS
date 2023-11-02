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
import {IconCheck, IconLoading, IconMinus, IconMinusCircle, IconPlus} from "@arco-design/web-react/icon";
import active from "@antv/g2/src/interaction/action/element/active";
import GetAxios from "@/utils/getaxios";
import UserAddonCountInfo from "@/pages/user/user-addon-count-info";

const { Title } = Typography;
const { Row, Col } = Grid;

const defaultVideoList = new Array(0).fill({});
export default function ListSearchResult() {
  const t = useLocale(locale);
  const [loading, setLoading] = useState(true);
  const [followLoading, setFollowLoading] = useState(false);
  const [followHovering, setFollowHovering] = useState(false);
  const [videoData, setVideoData] = useState(defaultVideoList);
  const [videoNum, setVideoNum] = useState({});
  const [videoNumU, setVideoNumU] = useState(0);

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

  // deep clone {}
  const cloneDeep = (obj) => {
    let newObj = {};
    for (let key in obj) {
      if (typeof obj[key] === 'object') {
        newObj[key] = cloneDeep(obj[key]);
      } else {
        newObj[key] = obj[key];
      }
    }
    return newObj;
  }

  const getVideoData = async (userid, t) => {
    setIsEndData(false);
    setLoading(true);
    let param = new FormData();
    param.append('user_id', isSelf ? (t === 'watched' ? 0 : userid ) : userid);
    param.append('relation', t);
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
          let tmp = cloneDeep(videoNum);
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
    param.append('relation', t);
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
          let tmp = cloneDeep(videoNum);
          tmp[t] = data.data.video_num;
          setVideoNum(tmp);
        }
      })
      .catch(error => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  }

  // todo: need to fix first enter page will not update the videoNum bug

  useEffect(() => {
    if (router.isReady && username) {
      setUserData(null);
      setVideoData(defaultVideoList);
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

  const followUser = (follow) => {
    setFollowLoading(true);
    const baxios = GetAxios();
    let params = new FormData();
    params.append('user_id', userData.user_id);
    params.append('action', follow?'unfollow':'follow');
    baxios.post('/v1-api/v1/user/follow', params)
      .then(response => {
        const data = response.data
        if (data.status !== 200) {
          console.error(data.err_msg);
          return;
        }
        (isSelf ? getSelfInfoData() : getUserInfoData(username))

      })
      .catch(error => {
        console.error(error);
      })
      .finally(() => {setFollowLoading(false);});
  }

  return (
    <div className={styles['container']}>
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
          <div className={styles['user-info']}>
            { /* add avatar to the left */}
            <Avatar size={64} style={{ }}>
              {userData?(userData.avatar_url?<img src={userData.avatar_url} />:userData.nickname):'A'}
            </Avatar>
            <div className={styles['base-info']}>
              <div className={styles.nickname}>
                {userData?(
                  <div>{userData.nickname}</div>
                ):''}
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
                    <Button type="outline" className={styles['change-nickname']}>修改昵称</Button>
                  </Popconfirm>
                ):userData?(
                  <Button
                    type={userData.be_followed ? "secondary" : "primary"}
                          className={styles['follow']}
                    icon={userData.be_followed ? (followHovering ? <IconMinusCircle /> : <IconCheck />) : <IconPlus />}
                    onClick={() => {followUser(userData.be_followed)}}
                    loading={followLoading}
                    onMouseEnter={() => setFollowHovering(true)}
                    onMouseLeave={() => setFollowHovering(false)}
                  >
                    {userData.be_followed ? (followHovering ? "取消" : "已") : null}关注
                  </Button>
                ):null}
              </div>
              <div className={styles.username}>@{userData?userData.username:''}</div>
            </div>
            <div className={styles['addon-info']}>
              <UserAddonCountInfo type={"关注"} data={userData?userData.follow_count:0} />
              <Divider type='vertical' style={{ height: "2em" }} />
              <UserAddonCountInfo type={"粉丝"} data={userData?userData.be_followed_count:0} />
              <Divider type='vertical' style={{ height: "2em" }} />
              <UserAddonCountInfo type={"获赞"} data={userData?userData.be_liked_count:0} />
              <Divider type='vertical' style={{ height: "2em" }} />
              <UserAddonCountInfo type={"浏览量"} data={userData?userData.be_watched_count:0} />
            </div>
          </div>
        )}
      </Card>
      { noSuchUser ? null : (
        <>
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
              {/*<Tabs.TabPane key="uploaded" title={t['cardList.tab.title.uploaded'] + (videoNumU!==0?" ("+videoNumU+")":"")} />*/}
              <Tabs.TabPane key="liked" title={t['cardList.tab.title.liked'] + (videoNum['liked']&&videoNum['liked']!==0?" ("+videoNum['liked']+")":"")} />
              <Tabs.TabPane key="favorite" title={t['cardList.tab.title.favorite'] + (videoNum['favorite']&&videoNum['favorite']!==0?" ("+videoNum['favorite']+")":"")} />
              {/* watched is only for self */ }
              {isSelf ? <Tabs.TabPane key="watched" title={t['cardList.tab.title.watched'] + (videoNum['watched']&&videoNum['watched']!==0?" ("+videoNum['watched']+")":"")} /> : null}
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
              { /* uploaded, liked, favorite, watched info */ }
              {`${isSelf ? '您' : '该用户'}还没有${activeKey === 'uploaded' ? '上传' : activeKey === 'liked' ? '点赞' : activeKey === 'favorite' ? '收藏' : '观看'}过任何视频`}
            </span>
                  </ContentContainer>
                ) : (
                  <ContentContainer>
                    <Button type='text' onClick={() => { getMoreData(userData.user_id, activeKey) }}>
                      {`加载更多${activeKey === 'uploaded' ? '上传' : activeKey === 'liked' ? '点赞' : activeKey === 'favorite' ? '收藏' : '观看过'}的视频`}
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
    </div>
  );
}
