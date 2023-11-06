import React, { useState, useEffect } from 'react';
import styles from './style/index.module.less';
import { GlobalState } from '@/store';
import {
  Button,
  Card,
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
import axios from 'axios';
import { IconCheck } from '@arco-design/web-react/icon';
import { useRouter } from 'next/router';
import { useSelector } from 'react-redux';
import baxios from "@/utils/getaxios";
import Head from "next/head";

const { Title, Paragraph } = Typography;

function UploadShortVideo() {
  const [loading, setLoading] = useState(false);
  const [current, setCurrent] = useState(1);
  const [fileList, setFileList] = useState([]);
  const [uploading, setUploading] = useState(false);
  const [videoUid, setVideoUid] = useState('');
  const [form] = Form.useForm();
  const [cancelTokenSource, setCancelTokenSource] = useState(null);

  const { isLogin, userLoading } = useSelector((state: GlobalState) => state);

  const router = useRouter();

  useEffect(() => {
    if (router.isReady) {
      if (userLoading !== undefined && !userLoading && !isLogin) {
        Message.error('请先登录');
        // window.location.href = '/';
        return;
      }
    }
  }, [router.isReady]);

  const onUploadChange = (files) => {
    setVideoUid('');
    setUploading(true);
    form.setFieldValue('video_uid', '');

    const newFiles = files.map((item) => ({
      ...item,
      percent: 0,
      status: 'uploading',
    }));
    setFileList(newFiles);

    // console.log(files);
    // refuse upload if file size is larger than 200MB
    if (files.length > 0 && files[0].originFile.size > 200 * 1024 * 1024) {
      console.log('file too big');
      // make file status to error
      setFileList((currentFileList) =>
        currentFileList.map((file) => {
          if (file.name === files[0].name) {
            return {
              ...file,
              percent: 0,
              status: 'error',
              response: '请选择小于200MB的短视频进行上传',
            };
          }
          return file;
        })
      );
      return;
    }

    if (files.length === 0) {
      return;
    }

    const formData = new FormData();
    formData.append('file', files[0].originFile);

    const source = axios.CancelToken.source();
    setCancelTokenSource(source);
    baxios
      .post('/v1-api/v1/video/upload/file', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        cancelToken: source.token,
        onUploadProgress: (progressEvent) => {
          const percentCompleted = Math.round(
            (progressEvent.loaded * 100) / progressEvent.total
          );
          setFileList((currentFileList) =>
            currentFileList.map((file) => {
              if (file.name === files[0].name) {
                return { ...file, percent: percentCompleted };
              }
              return file;
            })
          );
        },
      })
      .then((response) => {
        const data = response.data;
        console.log(data);
        if (data.status !== 200) {
          setFileList((currentFileList) =>
            currentFileList.map((file) => {
              if (file.name === files[0].name) {
                return { ...file, percent: 0, status: 'error' };
              }
              return file;
            })
          );
          console.error(data.err_msg);
          return;
        }
        // upload success, log video_uid
        const videoUid = data.data.video_uid;
        setVideoUid(videoUid);
        form.setFieldValue('video_uid', videoUid);

        // get content and keyword from filename which may like "aaaaa  ###k1 #k2 #k3 bbb  ##k1 #k2 ccc .mov"
        // content should be "aaaaa bbb ccc", keyword[] should be ['k1', 'k2', 'k3'] without duplicated
        // let's first extract keyword
        let keyword = [];
        let content = '';
        const filenameSplit = files[0].name.split('.');
        const contentAndKeyword = filenameSplit.slice(0, -1).join('.');
        // get all keyword after #
        const keywordSplit = contentAndKeyword.split('#');
        // get keywordSplit[1:] and remove empty string, and remove duplicated, and remove word after space
        // (keyword must be not split by space)
        keyword = keywordSplit
          .slice(1)
          .map((keyword) => {
            return keyword.trim().split(' ')[0];
          })
          .filter((keyword) => {
            return keyword !== '';
          })
          .filter((keyword, index, self) => {
            return self.indexOf(keyword) === index;
          });
        // remove keyword in contentAndKeyword to get content
        content = contentAndKeyword;
        keyword.forEach((keyword) => {
          content = content.replaceAll(keyword, '').replaceAll('#', '');
        });

        // replace possible multiple space such as "  " to " "
        content = content.replace(/\s+/g, ' ').trim();

        // add "#" back to keyword
        keyword = keyword.map((keyword) => {
          return '#' + keyword;
        });

        // console.log(content, keyword);

        // put content in filename to content if content is empty
        if (
          form.getFieldValue('content') === undefined ||
          form.getFieldValue('content') === ''
        ) {
          form.setFieldValue('content', content);
        }
        // put keyword in filename to keyword if keyword is empty
        if (
          form.getFieldValue('keyword') === undefined ||
          form.getFieldValue('keyword').length === 0
        ) {
          form.setFieldValue('keyword', keyword);
        }

        setFileList((currentFileList) =>
          currentFileList.map((file) => {
            if (file.name === files[0].name) {
              return { ...file, percent: 100, status: 'done' };
            }
            return file;
          })
        );
      })
      .catch((error) => {
        setFileList((currentFileList) =>
          currentFileList.map((file) => {
            if (file.name === files[0].name) {
              return { ...file, percent: 0, status: 'error' };
            }
            return file;
          })
        );
        console.error('Upload failed: ', error);
      })
      .finally(() => setUploading(false));
  };

  function submit() {
    console.log(form.getFields());
    console.log(videoUid);
    setLoading(true);

    const param = new FormData();
    param.append('video_uid', videoUid);
    param.append('video_type', form.getFieldValue('type'));
    param.append('video_content', form.getFieldValue('content'));
    param.append('video_keyword', form.getFieldValue('keyword').join(' '));
    baxios
      .post('/v1-api/v1/video/upload/confirm', param)
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

  const reCreateForm = () => {
    form.resetFields();
    setFileList([]);
    setVideoUid('');
    setCurrent(1);
  };

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

  const isAcceptFile = (file, accept) => {
    if (accept && file) {
      const accepts = Array.isArray(accept)
        ? accept
        : accept
            .split(',')
            .map((x) => x.trim())
            .filter((x) => x);
      const fileExtension =
        file.name.indexOf('.') > -1 ? file.name.split('.').pop() : '';
      return accepts.some((type) => {
        const text = type && type.toLowerCase();
        const fileType = (file.type || '').toLowerCase();
        if (text === fileType) {
          return true;
        }
        if (new RegExp('/*').test(text)) {
          const regExp = new RegExp('/.*$');
          return fileType.replace(regExp, '') === text.replace(regExp, '');
        }
        if (new RegExp('..*').test(text)) {
          return text === `.${fileExtension && fileExtension.toLowerCase()}`;
        }
        return false;
      });
    }
    return !!file;
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
    <>
      <Head>
        <title>短视频上传 - ACES短视频</title>
      </Head>
      <div className={styles.container}>
        <Card>
          <div className={styles.wrapper}>
            <Form form={form} className={styles.form}>
              {current === 1 && (
                <Form.Item noStyle>
                  <Title heading={4} style={{ marginBottom: '48px' }}>
                    {'填写信息上传短视频'}
                  </Title>
                  <Form.Item
                    label={'短视频文件'}
                    required
                    // requiredSymbol={false}
                    field="file"
                    validateTrigger={[]}
                    rules={[
                      {
                        validator: (value, callback) => {
                          // do not validate when uploading
                          if (uploading) {
                            callback('视频未上传完成');
                          } else if (videoUid === '') {
                            callback('请先上传视频文件');
                          } else {
                            callback();
                          }
                        },
                      },
                    ]}
                    extra={
                      '文件名中带有#的会被识别为关键词，除此外的其他内容会被识别为视频简介'
                    }
                  >
                    <div className={styles.customUpload}>
                      <Upload
                        disabled={!isLogin}
                        limit={1}
                        drag
                        accept="video/*"
                        showUploadList={{
                          startIcon: (
                            <Button size="mini" type="text">
                              开始上传
                            </Button>
                          ),
                          cancelIcon: (
                            <Button
                              size="mini"
                              type="text"
                              onClick={() => {
                                if (cancelTokenSource) {
                                  cancelTokenSource.cancel('Upload canceled.');
                                }
                                setFileList([]);
                                setVideoUid('');
                                form.setFieldValue('video_uid', '');
                              }}
                            >
                              取消上传
                            </Button>
                          ),
                          reuploadIcon: (
                            <Button size="mini" type="text">
                              点击重试
                            </Button>
                          ),
                          successIcon: (
                            <div>
                              <IconCheck />
                              <Button
                                size="mini"
                                type="text"
                                onClick={() => {
                                  setFileList([]);
                                  setVideoUid('');
                                  form.setFieldValue('video_uid', '');
                                }}
                              >
                                删除文件
                              </Button>
                            </div>
                          ),
                        }}
                        progressProps={{
                          size: 'default',
                          type: 'line',
                          showText: true,
                          width: '100%',
                        }}
                        // onDrop={(e) => {
                        //   let uploadFile = e.dataTransfer.files[0]
                        //   if (isAcceptFile(uploadFile, 'video/*')) {
                        //     return
                        //   } else {
                        //     // show error message
                        //   }
                        // }}
                        autoUpload={false}
                        onChange={onUploadChange}
                        fileList={fileList}
                      ></Upload>
                    </div>
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
                    <Select disabled={!isLogin}>
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
                      disabled={!isLogin}
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
                    <InputTag allowClear dragToSort disabled={!isLogin} />
                  </Form.Item>
                </Form.Item>
              )}
              {/* current === 2, confirm form items type, content, keyword on current === 1 */}
              {current === 2 && (
                <Form.Item noStyle>
                  <Title heading={4} style={{ marginBottom: '48px' }}>
                    {'确认短视频及信息'}
                  </Title>
                  <Form.Item label={'视频预览'}>
                    <video
                      src={URL.createObjectURL(fileList[0].originFile)}
                      autoPlay={true}
                      muted={true}
                      controls={true}
                      style={{
                        // width: '100%',
                        maxWidth: '494px',
                        maxHeight: '400px',
                      }}
                      onLoadedMetadata={(e) => {
                        const time = e.currentTarget.duration;
                        form.setFieldValue('time', time);
                      }}
                    />
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
                        disabled={loading || !isLogin}
                      >
                        返回修改
                      </Button>
                    )}
                    {current === 1 && (
                      <Button disabled={!isLogin} type="primary" size="large" onClick={toNext}>
                        下一步
                      </Button>
                    )}
                    {current === 2 && (
                      <Button
                        type="primary"
                        size="large"
                        loading={loading || !isLogin}
                        onClick={toConfirm}
                      >
                        {loading ? '发布中' : '确认发布'}
                      </Button>
                    )}
                  </Space>
                </Form.Item>
              ) : (
                <Form.Item noStyle>
                  <Result
                    status="success"
                    title={'提交成功'}
                    subTitle={'短视频上传成功！请等待3-5分钟转码'}
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
                      <Button key="again" type="primary" onClick={reCreateForm}>
                        {'再上传一个'}
                      </Button>,
                    ]}
                  />
                </Form.Item>
              )}
            </Form>
          </div>
          {current === 3 && (
            <div className={styles['form-extra']}>
              <Title heading={6}>{'视频上传说明'}</Title>
              <Paragraph type="secondary">
                {
                  '视频在上传成功后，将会在服务器进行转码以适应短视频播放，请在3-5分钟后再查看发布的视频。'
                }
              </Paragraph>
            </div>
          )}
        </Card>
      </div>
    </>
  );
}

export default UploadShortVideo;
