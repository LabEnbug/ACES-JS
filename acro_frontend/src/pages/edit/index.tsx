import React, { useState, useRef, useEffect } from 'react';
import styles from './style/index.module.less';
import { GlobalState } from '@/store';
import {
  Button,
  Card,
  DatePicker,
  Descriptions,
  Form,
  Input,
  InputTag,
  Message,
  Result,
  Select,
  Space,
  Tag,
  Typography,
  Upload,
} from '@arco-design/web-react';
import GetAxios from '@/utils/getaxios';
import axios, { Canceler } from 'axios';
import { IconCheck } from '@arco-design/web-react/icon';
import { useRouter } from 'next/router';
import MessageBox from '@/components/MessageBox';
import { useSelector } from 'react-redux';

const { Title, Paragraph } = Typography;

function UploadShortVideo() {
  const [loading, setLoading] = useState(false);
  const [current, setCurrent] = useState(1);
  const [uploading, setUploading] = useState(false);
  const [videoUid, setVideoUid] = useState('');
  const [form] = Form.useForm();
  const [cancelTokenSource, setCancelTokenSource] = useState(null);

  const [videoInfo, setVideoInfo] = useState(null);

  const { userInfo } = useSelector((state: GlobalState) => state);

  const router = useRouter();
  const { video_uid } = router.query;

  useEffect(() => {
    if (!video_uid) {
      window.location.href = '/';
    } else if (router.isReady && video_uid) {
      setVideoUid(video_uid.toString());
      // get video info
      const baxios = GetAxios();
      const param = new FormData();
      param.append('video_uid', video_uid.toString());
      baxios
        .post('/v1-api/v1/video/info', param)
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
            keyword: video.keyword.split(' '),
          });
          setVideoInfo(video);
        })
        .catch((error) => {
          console.error(error);
        })
        .finally(() => setLoading(false));
    }
  }, [router.isReady, video_uid]);

  function submit() {
    console.log(form.getFields());
    console.log(videoUid);
    setLoading(true);

    const baxios = GetAxios();
    const param = new FormData();
    param.append('video_uid', videoUid);
    param.append('video_type', form.getFieldValue('type'));
    param.append('video_content', form.getFieldValue('content'));
    param.append('video_keyword', form.getFieldValue('keyword').join(' '));
    baxios
      .post('/v1-api/v1/video/info/set', param)
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          Message.error(data.err_msg);
          return;
        }
        setCurrent(3);
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

  return (
    <div className={styles.container}>
      <Card>
        <div className={styles.wrapper}>
          <Form form={form} className={styles.form}>
            {current === 1 && (
              <Form.Item noStyle>
                <Title heading={4} style={{ marginBottom: '48px' }}>
                  {'修改短视频信息'}
                </Title>
                <Form.Item label={'短视频封面'}>
                  {videoInfo && videoInfo.cover_url !== '' && (
                    <img
                      style={{
                        width: '100%',
                        maxWidth: '494px',
                        maxHeight: '400px',
                      }}
                      src={videoInfo.cover_url}
                    ></img>
                  )}
                </Form.Item>
                <Form.Item
                  label={'视频类型'}
                  initialValue="4"
                  field="type"
                  rules={[
                    {
                      required: true,
                      message: '请选择视频类型',
                    },
                  ]}
                >
                  <Select>
                    {Object.keys(typeMap).map((key) => (
                      <Select.Option key={key} value={key}>
                        {typeMap[key]}
                      </Select.Option>
                    ))}
                  </Select>
                </Form.Item>
                <Form.Item
                  label={'视频简介'}
                  field="content"
                  defaultValue={''}
                  rules={[
                    {
                      required: true,
                      message: '请至少输入一个字符',
                    },
                    {
                      maxLength: 120,
                      message: '请将视频简介限制在120字内',
                    },
                  ]}
                >
                  <Input.TextArea
                    maxLength={{ length: 120, errorOnly: true }}
                    showWordLimit
                    autoSize={{ minRows: 1, maxRows: 6 }}
                  />
                </Form.Item>
                <Form.Item
                  label={'关键词'}
                  initialValue={[]}
                  field="keyword"
                  extra={'输入后回车生成'}
                  // add # before each keyword, but do not duplicate
                  normalize={(value) => {
                    return value
                      .map((keyword) => {
                        keyword = keyword.trim();
                        //delete multiple # at the beginning
                        while (keyword.startsWith('#')) {
                          keyword = keyword.slice(1);
                        }
                        return '#' + keyword;
                      })
                      .filter((keyword, index, self) => {
                        return self.indexOf(keyword) === index;
                      });
                  }}
                >
                  <InputTag allowClear dragToSort />
                </Form.Item>
              </Form.Item>
            )}
            {/* current === 2, confirm form items type, content, keyword on current === 1 */}
            {current === 2 && (
              <Form.Item noStyle>
                <Title heading={4} style={{ marginBottom: '48px' }}>
                  {'确认信息'}
                </Title>
                <Form.Item label={'短视频封面'}>
                  {videoInfo && videoInfo.cover_url !== '' && (
                    <img
                      style={{
                        width: '100%',
                        maxWidth: '494px',
                        maxHeight: '400px',
                      }}
                      src={videoInfo.cover_url}
                    ></img>
                  )}
                </Form.Item>
                <Form.Item label={'视频类型'}>
                  <Typography.Text>
                    {typeMap[form.getFieldValue('type')]}
                  </Typography.Text>
                </Form.Item>
                <Form.Item label={'视频简介'}>
                  <Typography.Text>
                    {form.getFieldValue('content')}
                  </Typography.Text>
                </Form.Item>
                <Form.Item label={'关键词'}>
                  <Typography.Text>
                    {form.getFieldValue('keyword').map((keyword, index) => (
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
                    ))}
                  </Typography.Text>
                </Form.Item>
              </Form.Item>
            )}
            {current !== 3 ? (
              <Form.Item label=" ">
                <Space>
                  {current === 2 && (
                    <Button
                      size="large"
                      onClick={() => setCurrent(current - 1)}
                      disabled={loading}
                    >
                      返回修改
                    </Button>
                  )}
                  {current === 1 && (
                    <Button type="primary" size="large" onClick={toNext}>
                      下一步
                    </Button>
                  )}
                  {current === 2 && (
                    <Button
                      type="primary"
                      size="large"
                      loading={loading}
                      onClick={toConfirm}
                    >
                      {loading ? '修改中' : '确认修改'}
                    </Button>
                  )}
                </Space>
              </Form.Item>
            ) : (
              <Form.Item noStyle>
                <Result
                  status="success"
                  title={'提交成功'}
                  subTitle={'短视频信息修改成功！'}
                  extra={[
                    <Button
                      key="watch"
                      style={{ marginRight: 16 }}
                      onClick={() => {
                        router.push({
                          pathname: `/video`,
                          query: {
                            video_uid: videoUid,
                          },
                        });
                      }}
                    >
                      {'前往观看该短视频'}
                    </Button>,
                    <Button
                      key="watchPage"
                      style={{ marginRight: 16 }}
                      onClick={() => {
                        router.push({
                          pathname: `/user/self`,
                        });
                      }}
                    >
                      {'查看已上传的短视频'}
                    </Button>,
                  ]}
                />
              </Form.Item>
            )}
          </Form>
        </div>
      </Card>
    </div>
  );
}

export default UploadShortVideo;
