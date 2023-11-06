import React, { useState, useEffect } from 'react';
import styles from './style/index.module.less';
import { GlobalState } from '@/store';
import {
  Button,
  Card, Divider,
  Form,
  Input, InputNumber,
  Message,
  Result,
  Space, Tag,
  Typography,
} from '@arco-design/web-react';
import { useRouter } from 'next/router';
import {useDispatch, useSelector} from 'react-redux';
import baxios from "@/utils/getaxios";
import {UpdateUserInfoOnly} from "@/utils/getuserinfo";
import Head from "next/head";
import useLocale from "@/utils/useLocale";
import locale from "./locale"

const { Title } = Typography;

function Promote() {
  const t = useLocale(locale);
  const tg = useLocale();
  const [loading, setLoading] = useState(false);
  const [current, setCurrent] = useState(1);
  const [videoUid, setVideoUid] = useState('');
  const [form] = Form.useForm();
  const [videoInfo, setVideoInfo] = useState(null);
  const [remainPromoteCount, setRemainPromoteCount] = useState(0);
  const [price, setPrice] = useState(0.10);


  const { isLogin, userLoading, userInfo, init } = useSelector((state: GlobalState) => state);
  const [isUserUploaded, setIsUserUploaded] = useState(false);

  const router = useRouter();
  const { video_uid } = router.query;

  const dispatch = useDispatch();

  const perPrice = 0.1; // Promotion per visit price
  const maxCount = 10000; // Max promotion count

  function GetVideoInfo() {
    baxios
      .get('/videos/' + video_uid.toString())
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          Message.error(data.err_msg);
          window.location.href = '/';
          return;
        }
        const video = data.data.video;
        form.setFieldsValue({
          type: video.type.toString(),
          content: video.content,
          keyword: video.keyword.split(' ').filter((item) => item !== ''),
        });
        setVideoInfo(video);
        setRemainPromoteCount(data.data.remain_promote_count)
        setIsUserUploaded(video.is_user_uploaded)
        if (!video.is_user_uploaded) {
          Message.error('暂未开放推广其他用户的视频');
        }
      })
      .catch((error) => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  }

  useEffect(() => {
    if (router.isReady) {
      if (init && userLoading!==undefined && !userLoading && !isLogin) {
        Message.error('请先登录');
        // window.location.href = '/';
        return;
      }
    }

    if (!video_uid) {
      Message.error('无此视频');
      // window.location.href = '/';
    } else if (router.isReady && video_uid) {
      setVideoUid(video_uid.toString());
      // get video info
      GetVideoInfo();
    }
  }, [router.isReady, video_uid]);

  function submit() {
    console.log(form.getFields());
    if (!isBalanceEnough()) {
      Message.error('余额不足，请先充值。');
      return;
    }
    setLoading(true);

    const param = new FormData();
    param.append('count', form.getFieldValue('count'));
    baxios
      .post('/videos/' + videoUid + '/actions/' + 'promote', param)
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          Message.error(data.err_msg);
          return;
        }

        UpdateUserInfoOnly(dispatch);
        GetVideoInfo();
        setCurrent(2);
      })
      .catch((error) => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  }

  const typeMap = {
    1: '知识',
    2: '热点',
    3: '游戏',
    4: '娱乐',
    5: '二次元',
    6: '音乐',
    7: '美食',
    8: '体育',
    9: '时尚',
  };

  const toNext = async () => {
    try {
      await form.validate();
      setCurrent(current + 1);
    } catch (_) {}
  };

  const toConfirm = () => {
    submit();
  };

  const isBalanceEnough = () => {
    return price <= userInfo.balance;
  }

  const onValuesChange = (changedValues: any, allValues: any) => {
    if (changedValues.count) {
      setPrice(changedValues.count * perPrice);
    }
  }

  return (
    <>
      <Head>
        <title>{t['title']} - {tg['title.global']}</title>
      </Head>
      <div className={styles.container}>
        <Card>
          <div className={styles.wrapper}>
            <Form form={form} className={styles.form} onValuesChange={onValuesChange}>
              {current === 1 && (
                <Form.Item noStyle>
                  <Title heading={4} style={{ marginBottom: '48px' }}>
                    {'短视频推广'}
                  </Title>
                  <Form.Item label={'短视频封面'}>
                    {videoInfo && videoInfo.cover_url !== '' && (
                      <img
                        style={{
                          width: '100%',
                          height: '100%',
                          objectFit: 'contain',
                          maxWidth: '494px',
                          maxHeight: '400px',
                        }}
                        src={videoInfo.cover_url}
                      ></img>
                    )}
                  </Form.Item>
                  <Form.Item label={'视频类型'}>
                    <Typography.Text>
                      {typeMap[videoInfo?videoInfo['type']:1]}
                    </Typography.Text>
                  </Form.Item>
                  <Form.Item label={'视频简介'}>
                    <Typography.Text>
                      {videoInfo?videoInfo['content']:''}
                    </Typography.Text>
                  </Form.Item>
                  <Form.Item label={'关键词'}>
                    <Typography.Text>
                      {videoInfo ? (videoInfo['keyword'].split(' ').filter((item) => item !== '').map((keyword, index) => (
                        <Tag
                          key={index.toString()}
                          style={{
                            cursor: 'pointer',
                            marginRight: '4px',
                            marginBottom: '4px',
                            backgroundColor: 'rgba(var(--gray-6), 0.3)',
                          }}
                        >
                          {keyword}
                        </Tag>
                      ))) : null}
                    </Typography.Text>
                  </Form.Item>
                  <Divider ></Divider>
                  <Form.Item label={'剩余推广数'}>
                    <Typography.Text>
                      {isLogin&&userInfo?remainPromoteCount + ' 次':''}
                    </Typography.Text>
                  </Form.Item>
                  <Form.Item label={'账户余额'}>
                    <Typography.Text>
                      {isLogin&&userInfo?(userInfo.balance).toFixed(2) + ' 元':''}
                    </Typography.Text>
                  </Form.Item>
                  <Form.Item
                    label={'流量数'}
                    initialValue="1"
                    field="count"
                    extra={'要推广的数量，当被推广用户播放一次该视频，则扣除一次流量数'}
                    rules={[
                      {
                        required: true,
                        message: '请填写流量数',
                      },
                      {
                        validateTrigger: 'onBlur',
                        validator: (v, cb) => {
                          if (!isBalanceEnough()) {
                            return cb('余额不足');
                          }
                          return null;
                        },
                      }
                    ]}
                  >
                    <InputNumber autoComplete={'off'} max={maxCount} min={1} disabled={!isLogin || !isUserUploaded}/>
                  </Form.Item>
                  <Form.Item label={'需扣除余额'}>
                    <Typography.Text>
                      {price.toFixed(2)} 元
                      {(!isBalanceEnough() && <div style={{color: 'orange'}}>(余额不足, 仍需充值 {(price-userInfo.balance).toFixed(2)} 元)</div>)}
                    </Typography.Text>
                  </Form.Item>
                </Form.Item>
              )}
              {current === 1 && (
                <Form.Item label=" ">
                  <Space>
                    {current === 1 && (
                        <Button
                        type="primary"
                        size="large"
                        disabled={loading || !isLogin || !isUserUploaded}
                        loading={loading}
                        onClick={toConfirm}
                      >
                        {loading ? '提交中' : '确认推广'}
                      </Button>
                    )}
                  </Space>
                </Form.Item>
              )}
              {current === 2 && (
                <Form.Item noStyle>
                  <Result
                    status="success"
                    title={'推广成功'}
                    subTitle={'成功为视频充值推广流量 ' + form.getFieldValue("count") + ' 次'}
                    extra={[
                      <Button
                        key="watchPage"
                        style={{ marginRight: 16 }}
                        onClick={() => {
                          router.push({
                            pathname: `/user/self`,
                          });
                        }}
                      >
                        {'进入用户页'}
                      </Button>,
                      <Button key="again" type="primary" onClick={() => {setCurrent(current - 1)}}>
                        {'再推广该视频'}
                      </Button>,
                    ]}
                  />
                </Form.Item>
              )}
            </Form>
          </div>
        </Card>
      </div>
    </>
  );
}

export default Promote;
