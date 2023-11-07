import React, { useEffect, useRef, useState } from 'react';
import {
  Tabs,
  Card,
  Input,
  Button,
  List,
  Divider,
  Avatar, Message, Upload,
} from '@arco-design/web-react';
import useLocale from '@/utils/useLocale';
import locale from './locale';
import styles from './style/index.module.less';
import CardBlock from './card-block';
import { useRouter } from 'next/router';
import { Empty } from '@arco-design/web-react';
import { Popconfirm } from '@arco-design/web-react';
import {
  IconCamera,
  IconCheck,
  IconMinusCircle,
  IconPlus,
} from '@arco-design/web-react/icon';
import UserAddonCountInfo from '@/pages/user/user-addon-count-info';
import baxios from "@/utils/getaxios";
import Head from "next/head";
import {UpdateUserInfoOnly} from "@/utils/getuserinfo";
import {useDispatch} from "react-redux";

const defaultVideoList = new Array(0).fill({});
export default function UserPage() {
  const t = useLocale(locale);
  const tg = useLocale();
  const [loading, setLoading] = useState(true);
  const [followLoading, setFollowLoading] = useState(false);
  const [followHovering, setFollowHovering] = useState(false);
  const [videoData, setVideoData] = useState(defaultVideoList);
  // const [videoNumList, setVideoNumList] = useState({});
  const [uploadedNum, setUploadedNum] = useState(0);
  const [likedNum, setLikedNum] = useState(0);
  const [favoriteNum, setFavoriteNum] = useState(0);
  const [watchedNum, setWatchedNum] = useState(0);

  const [userData, setUserData] = useState(null);
  const [nicknamePopupVisible, setNicknamePopupVisible] = useState(false);
  const [nicknameForChange, setNicknameForChange] = useState('');

  const [activeKey, setActiveKey] = useState('uploaded');

  const router = useRouter();
  const { username } = router.query;

  const listRef = useRef(null);

  const [isEndData, setIsEndData] = useState(false);

  const [isSelf, setIsSelf] = useState(false);

  const [noSuchUser, setNoSuchUser] = useState(false);

  const [avatarFile, setAvatarFile] = useState(null);

  const dispatch = useDispatch();

  const setNum = (relation: string, num: number) => {
    if (relation === 'uploaded') {
      setUploadedNum(num);
    } else if (relation === 'liked') {
      setLikedNum(num);
    } else if (relation === 'favorite') {
      setFavoriteNum(num);
    } else if (relation === 'watched') {
      setWatchedNum(num);
    }
  }

  const getNum = (relation: string) => {
    if (relation === 'uploaded') {
      return uploadedNum;
    } else if (relation === 'liked') {
      return likedNum;
    } else if (relation === 'favorite') {
      return favoriteNum;
    } else if (relation === 'watched') {
      return watchedNum;
    }
  }
  const getUserInfoData = async (username) => {
    setLoading(true);
    baxios
      .get('/users/' + username)
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          setNoSuchUser(true);
          return;
        }
        setIsSelf(data.data.user.is_self);
        setUserData(data.data.user);
        getVideoData(data.data.user.user_id, 'uploaded');
      })
      .catch((error) => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  };

  const getSelfInfoData = async () => {
    setLoading(true);
    baxios
      .get('/user/info')
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          window.location.href = '/'; // reject not logged in
          return;
        }
        setUserData(data.data.user);
        getVideoData(data.data.user.user_id, 'uploaded');
      })
      .catch((error) => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  };

  // deep clone {}
  const cloneDeep = (obj) => {
    const newObj = {};
    for (const key in obj) {
      if (typeof obj[key] === 'object') {
        newObj[key] = cloneDeep(obj[key]);
      } else {
        newObj[key] = obj[key];
      }
    }
    return newObj;
  };

  const getVideoData = (userid, relation) => {
    setIsEndData(false);
    setLoading(true);
    // sleep
    // await new Promise(resolve => setTimeout(resolve, 3000));
    baxios
      .get('/videos?' +
        'user_id=' + (isSelf ? (relation === 'watched' ? 0 : userid) : userid) + '&' +
        'relation=' + relation + '&' +
        'limit=12')
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          // setVideoData(defaultVideoList);
          setNum(relation, 0)
          return;
        }
        setVideoData(data.data.video_list);
        if (data.data.video_num) {
          setNum(relation, data.data.video_num)
        }
      })
      .catch((error) => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  };

  const getMoreData = (userid, relation) => {
    setLoading(true);
    baxios
      .get('/videos?' +
      'user_id=' + (isSelf ? (relation === 'watched' ? 0 : userid) : userid) + '&' +
        'relation=' + relation + '&' +
        'start=' + videoData.length.toString() + '&' +
        'limit=12')
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          setIsEndData(true);
          return;
        }
        setVideoData(videoData.concat(data.data.video_list));
        if (data.data.video_num) {
          setNum(relation, data.data.video_num)
        }
      })
      .catch((error) => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    if (router.isReady && username) {
      // need to fix when user page to other user page, videoNum would not change
      // 20231106 fixed
      setNum('uploaded', 0);
      setNum('liked', 0);
      setNum('favorite', 0);
      setNum('watched', 0);

      setUserData(null);
      setVideoData(defaultVideoList);
      setActiveKey('uploaded');
      username === 'self' ? setIsSelf(true) : setIsSelf(false);
      username === 'self' ? getSelfInfoData() : getUserInfoData(username);
    }
  }, [router.isReady, username]);

  useEffect(() => {
    // get video after user data set
    if (userData) {
      setVideoData(defaultVideoList);
      getVideoData(userData.user_id, activeKey);
    }
  }, [activeKey]);

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
    <div style={{ textAlign: 'center', marginTop: 14 }}>
      <span style={{ color: 'var(--color-text-3)' }}>
        {/*<IconLoading style={{marginRight: 8, color: 'rgb(var(--arcoblue-6))'}}/>*/}
        {t['card.loading']}
      </span>
    </div>
  );

  function handleChangeNickname() {
    return new Promise<void>((resolve, reject) => {
      const params = new FormData();
      params.append('type', 'nickname');
      params.append('nickname', nicknameForChange);
      baxios
        .put('/user/info', params)
        .then((response) => {
          const data = response.data;
          if (data.status !== 200) {
            console.error(data.err_msg);
            Message.error(data.err_msg);
            resolve()
            return;
          }
          setUserData(data.data.user);
          Message.success("昵称修改成功！");
          UpdateUserInfoOnly(dispatch);
          setNicknamePopupVisible(false);
          resolve();
        })
        .catch((error) => {
          console.error(error);
          reject();
        })
        .finally();
    });
  }

  const handleChangeAvatar = (currentFile) => {
    const params = new FormData();
    params.append('type', 'avatar');
    params.append('file', currentFile.originFile);
    baxios
      .put('/user/info', params)
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          Message.error(data.err_msg);
          return;
        }
        setUserData(data.data.user);
        Message.success("头像修改成功！");
        UpdateUserInfoOnly(dispatch);
        setAvatarFile({
          ...currentFile,
          url: URL.createObjectURL(currentFile.originFile),
        });
      })
      .catch((error) => {
        console.error(error);
      })
      .finally();
  }

  const followUser = (follow) => {
    setFollowLoading(true);
    // sleep 1000ms
    setTimeout(() => {
      (follow ? baxios.delete : baxios.post)
      ('/users/' + userData.username + '/follow')
        .then((response) => {
          const data = response.data;
          if (data.status !== 200) {
            console.error(data.err_msg);
            Message.error("请先登录");
            return;
          }
          isSelf ? getSelfInfoData() : getUserInfoData(username);
        })
        .catch((error) => {
          console.error(error);
        })
        .finally(() => {
          setFollowLoading(false);
        });
    }, 200);
  };

  const handleDelete = (itemToDelete) => {
    setVideoData(videoData.filter((item) => item !== itemToDelete));
    setNum(activeKey, getNum(activeKey) - 1);
  };

  return (
    <>
      <Head>
        <title>{userData?userData.nickname + ' - ':''}{t['title']} - {tg['title.global']}</title>
      </Head>
      <div className={styles['container']}>
        <Card className={styles['top-user-info-wrapper']}>
          {noSuchUser ? (
            <div style={{ textAlign: 'center', marginTop: 16 }}>
              <div
                style={{
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                }}
              >
                <span style={{ color: 'var(--color-text-1)' }}>{t['user.notExist']}</span>
              </div>
            </div>
          ) : (
            <div className={styles['top-user-info']}>
              {/* add avatar to the left */}
              <div className={styles['top-user-card-left']}>
                <Upload
                  showUploadList={false}
                  fileList={avatarFile ? [avatarFile] : []}
                  accept="image/*"
                  autoUpload={false}
                  onChange={(_, currentFile) => {
                    // refuse upload if file size is larger than 200MB
                    console.log(currentFile);
                    if (currentFile.originFile.size > 2 * 1024 * 1024) {
                      console.log('file too big');
                      Message.error("请选择小于2MB大小的头像图片进行上传");
                      return;
                    }
                    handleChangeAvatar(currentFile);
                  }}
                >
                  <Avatar
                    className={styles['avatar']}
                    size={64}
                    triggerIcon={isSelf?<IconCamera />:null}
                    triggerType='mask'
                  >
                    {avatarFile && avatarFile.url ? (
                      <img src={avatarFile.url}  alt={null}/>
                    ) : (
                      userData ? (
                      userData.avatar_url ? (
                        <img src={userData.avatar_url}  alt={null}/>
                      ) : (
                        userData.nickname
                      )
                    ) : ('A')
                      )}
                  </Avatar>
                </Upload>
              </div>
              <div className={styles['top-user-card-right']}>
                <div className={styles.nickname}>
                  {userData && <div>{userData.nickname}</div>}
                  {isSelf ? (
                    <Popconfirm
                      position="bottom"
                      icon={null}
                      popupVisible={nicknamePopupVisible}
                      title={
                        <Input
                          autoComplete={'off'}
                          placeholder={t['user.change.nickname.placeholder']}
                          value={nicknameForChange}
                          onChange={(e) => setNicknameForChange(e)}
                          onKeyPress={(e) => {
                            if (e.key === 'Enter') {
                              handleChangeNickname().then(() => {setNicknamePopupVisible(false)});
                            }
                          }}
                        />
                      }
                      okText={t['user.change.nickname.ok']}
                      cancelText={t['user.change.nickname.cancel']}
                      onOk={handleChangeNickname}
                      onCancel={() => {
                        setNicknamePopupVisible(false);
                      }}
                    >
                      <Button
                        type="outline"
                        className={styles['change-nickname']}
                        onClick={() => setNicknamePopupVisible(true)}
                      >
                        {t['user.change.nickname']}
                      </Button>
                    </Popconfirm>
                  ) : (
                    userData && (
                      <Button
                        type={userData.be_followed ? 'secondary' : 'primary'}
                        className={styles['follow']}
                        icon={
                          userData.be_followed ? (
                            followHovering ? (
                              <IconMinusCircle />
                            ) : (
                              <IconCheck />
                            )
                          ) : (
                            <IconPlus />
                          )
                        }
                        onClick={() => {
                          followUser(userData.be_followed);
                        }}
                        loading={followLoading}
                        onMouseEnter={() => setFollowHovering(true)}
                        onMouseLeave={() => setFollowHovering(false)}
                      >
                        {userData.be_followed ? (
                          (followHovering ? t['user.cancel'] : t['user.already']) + t['user.followed']
                        ) : (t['user.follow'])}
                      </Button>
                    )
                  )}
                </div>
                <div className={styles.username}>
                  @{userData ? userData.username : ''}
                </div>
              </div>
              <div className={styles['top-user-addon-info']}>
                <UserAddonCountInfo
                  type={t['user.followed']}
                  data={userData ? userData.follow_count : 0}
                />
                <Divider type="vertical" style={{ height: '2em' }} />
                <UserAddonCountInfo
                  type={t['user.follower']}
                  data={userData ? userData.be_followed_count : 0}
                />
                <Divider type="vertical" style={{ height: '2em' }} />
                <UserAddonCountInfo
                  type={t['user.liked']}
                  data={userData ? userData.be_liked_count : 0}
                />
                <Divider type="vertical" style={{ height: '2em' }} />
                <UserAddonCountInfo
                  type={t['user.view']}
                  data={userData ? userData.be_watched_count : 0}
                />
              </div>
            </div>
          )}
        </Card>
        {!noSuchUser && (
          <>
            <Card>
              {/*<Title heading={6}>{t['menu.list.card']}</Title>*/}
              <Tabs activeTab={activeKey} type="text" onChange={setActiveKey}>
                <Tabs.TabPane
                  key="uploaded"
                  title={
                    t['cardList.tab.title.uploaded'] +
                    (getNum('uploaded') !== 0 ? ' (' + getNum('uploaded') + ')' : '')
                  }
                />
                {/*<Tabs.TabPane key="uploaded" title={t['cardList.tab.title.uploaded'] + (videoNumU!==0?" ("+videoNumU+")":"")} />*/}
                <Tabs.TabPane
                  key="liked"
                  title={
                    t['cardList.tab.title.liked'] +
                    (getNum('liked') !== 0 ? ' (' + getNum('liked') + ')' : '')
                  }
                />
                <Tabs.TabPane
                  key="favorite"
                  title={
                    t['cardList.tab.title.favorite'] +
                    (getNum('favorite') !== 0 ? ' (' + getNum('favorite') + ')' : '')
                  }
                />
                {/* watched is only for self */}
                {isSelf ? (
                  <Tabs.TabPane
                    key="watched"
                    title={
                      t['cardList.tab.title.watched'] +
                      (getNum('watched') !== 0 ? ' (' + getNum('watched') + ')' : '')
                    }
                  />
                ) : null}
              </Tabs>
              <Divider />
              <List
                ref={listRef}
                className={styles['card-list']}
                grid={{
                  xs: 12,
                  sm: 12,
                  md: 12,
                  lg: 8,
                  xl: 8,
                  xxl: 4,
                }}
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
                            }}
                          >
                            {/* uploaded, liked, favorite, watched info */}
                            {`${isSelf ? t['card.you'] : t['card.this.user']}${t['card.never']}${
                              activeKey === 'uploaded'
                                ? t['card.uploaded']
                                : activeKey === 'liked'
                                ? t['card.liked']
                                : activeKey === 'favorite'
                                ? t['card.favorite']
                                : t['card.watched']
                            }${t['card.any.video']}`}
                          </span>
                        </ContentContainer>
                      }
                    ></Empty>
                  )
                }
                style={{ overflowY: 'scroll', height: 'calc(100vh - 200px)' }}
                dataSource={videoData}
                bordered={false}
                onReachBottom={() => {
                  userData && getMoreData(userData.user_id, activeKey);
                }}
                render={(item, _) => (
                  <List.Item style={{ padding: '4px 4px' }}>
                    {
                      <CardBlock
                        card={item}
                        onDelete={() => handleDelete(item)}
                        type={activeKey}
                        watching_username={userData ? userData.username : ''}
                        loading={loading}
                      />
                    }
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
                        style={{
                          color: 'var(--color-text-3)',
                          marginBottom: '4px',
                        }}
                      >
                        {t['card.noMoreContent']}
                      </span>
                    </ContentContainer>
                  ) : (
                    videoData.length !== 0 && (
                      <ContentContainer>
                        <Button
                          type="text"
                          onClick={() => {
                            getMoreData(userData.user_id, activeKey);
                          }}
                        >
                          {`${t['card.loadMore']}${
                            activeKey === 'uploaded'
                              ? t['card.uploaded']
                              : activeKey === 'liked'
                              ? t['card.liked']
                              : activeKey === 'favorite'
                              ? t['card.favorite']
                              : t['card.watched']
                          }${t['card.video']}`}
                        </Button>
                      </ContentContainer>
                    )
                  )
                }
              />
            </Card>
          </>
        )}
      </div>
    </>
  );
}
