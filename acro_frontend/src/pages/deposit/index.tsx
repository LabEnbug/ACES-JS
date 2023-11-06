import React, { useState, useEffect } from 'react';
import styles from './style/index.module.less';
import { GlobalState } from '@/store';
import {
  Button,
  Card,
  Form,
  Input,
  Message,
  Result,
  Space,
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

function Deposit() {
  const t = useLocale(locale);
  const tg = useLocale();
  const [loading, setLoading] = useState(false);
  const [current, setCurrent] = useState(1);

  const [form] = Form.useForm();
  const [deposit_amount, setDepositAmount] = useState(0);

  const { isLogin, userLoading, userInfo, init } = useSelector((state: GlobalState) => state);

  const router = useRouter();
  const dispatch = useDispatch();

  useEffect(() => {
    if (router.isReady) {
      if (init && userLoading!==undefined && !userLoading && !isLogin) {
        Message.error('请先登录');
        // window.location.href = '/';
        return;
      }
    }
  }, []);

  function submit() {
    console.log(form.getFields());
    setLoading(true);

    const param = new FormData();
    param.append('card_key', form.getFieldValue('cardKey'));
    baxios
      .post('/user/deposit', param)
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          Message.error(data.err_msg);
          return;
        }
        setDepositAmount(data.data.deposit_amount);
        UpdateUserInfoOnly(dispatch);
        setCurrent(2);
      })
      .catch((error) => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  }

  const toNext = async () => {
    try {
      await form.validate();
      setCurrent(current + 1);
    } catch (_) {}
  };

  const toConfirm = () => {
    submit();
  };

  const onValuesChange = (changedValues: any, allValues: any) => {
    // console.log(changedValues, allValues);
    //set cardKey up upper
    if (changedValues.cardKey) {
      if (changedValues.cardKey.length < 5) {
        form.setFieldsValue({
          cardKey: 'ACES-',
        });
        return;
      }

      // make - between every 4 chars, split
      const arr = changedValues.cardKey.replaceAll('-', '');
      const newArr = [];
      for (let i = 0; i < arr.length; i += 4) {
        newArr.push(arr.slice(i, i + 4));
      }
      const newStr = newArr.join('-');

      const upper = newStr.toUpperCase();
      if (upper !== newStr) {
        form.setFieldsValue({
          cardKey: upper,
        });
      }
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
                    {'开始充值'}
                  </Title>
                  <Form.Item label={'账户余额'}>
                    <Typography.Text>
                      {isLogin&&userInfo?userInfo.balance + ' 元':''}
                    </Typography.Text>
                  </Form.Item>
                  <Form.Item
                    label={'充值卡号'}
                    initialValue="ACES-"
                    field="cardKey"
                    extra={'格式为：ACES-XXXX-XXXX-XXXX'}
                    rules={[
                      {
                        required: true,
                        message: '请填写充值卡号',
                      },
                      {
                        validateTrigger: 'onBlur',
                        match: /^ACES-.+$/,
                        message: '请填写充值卡号',
                      },
                      {
                        validateTrigger: 'onBlur',
                        match: /^ACES-([0-9A-Z]{4}-){2}[0-9A-Z]{4}$/,
                        message: '充值卡号格式错误',
                      },
                    ]}
                  >
                    <Input
                      autoComplete={'off'}
                      disabled={!isLogin}
                      maxLength={{ length: 19 }}
                      showWordLimit
                    />
                  </Form.Item>
                </Form.Item>
              )}
              {current !== 2 ? (
                <Form.Item label=" ">
                  <Space>
                    {/*{current === 2 && (*/}
                    {/*  <Button*/}
                    {/*    size="large"*/}
                    {/*    onClick={() => setCurrent(current - 1)}*/}
                    {/*    disabled={loading || !isLogin}*/}
                    {/*  >*/}
                    {/*    返回修改*/}
                    {/*  </Button>*/}
                    {/*)}*/}
                    {current === 1 && (
                        <Button
                        type="primary"
                        size="large"
                        disabled={loading || !isLogin}
                        loading={loading}
                        onClick={toConfirm}
                      >
                        {loading ? '充值中' : '确认充值'}
                      </Button>
                    )}
                    {/*{current === 2 && (*/}
                    {/*  <Button*/}
                    {/*    type="primary"*/}
                    {/*    size="large"*/}
                    {/*    loading={loading || !isLogin}*/}
                    {/*    onClick={toConfirm}*/}
                    {/*  >*/}
                    {/*    {loading ? '充值中' : '确认充值'}*/}
                    {/*  </Button>*/}
                    {/*)}*/}
                  </Space>
                </Form.Item>
              ) : (
                <Form.Item noStyle>
                  <Result
                    status="success"
                    title={'充值成功'}
                    subTitle={'成功为账户充值 ' + deposit_amount + ' 元'}
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
                      <Button key="again" type="primary" onClick={() => setCurrent(current - 1)}>
                        {'再充值一个'}
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

export default Deposit;
