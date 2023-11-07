import React, { useContext, useEffect, useState } from 'react';
import {
  Input,
  Avatar,
  Select,
  Dropdown,
  Menu,
  Divider,
  Message,
  Button,
  Modal,
  Form,
} from '@arco-design/web-react';
import {
  IconLanguage,
  IconUser,
  IconSettings,
  IconPoweroff,
  IconLoading,
  IconLock,
  IconUpload, IconPen,
} from '@arco-design/web-react/icon';
import IconPaymentCheck from '@/assets/payment-check.svg';
import { useSelector, useDispatch } from 'react-redux';
import { GlobalState } from '@/store';
import { GlobalContext } from '@/context';
import useLocale from '@/utils/useLocale';
import Logo from '@/assets/logo.svg';
import SearchPopupBox from '@/components/SearchPopupBox';
import IconButton from './IconButton';
import Settings from '../Settings';
import styles from './style/index.module.less';
import defaultLocale from '@/locale';
import { useRouter } from 'next/router';
import {setToken} from "@/utils/authentication";
import baxios from "@/utils/getaxios";
import FetchUserInfo from "@/utils/getuserinfo";
import Link from "next/link";

const FormItem = Form.Item;

function Navbar({ show }: { show: boolean }) {
  const t = useLocale();
  const { userInfo, userLoading, isLogin } = useSelector((state: GlobalState) => state);
  const dispatch = useDispatch();
  // const [_, setUserStatus] = useStorage('userStatus');
  // const [role, setRole] = useStorage('userRole', 'user');
  const [signInModal, SetSignInModal] = useState(false);
  const [signUpModal, SetSignUpModal] = useState(false);

  const router = useRouter();
  const [searchKeyword, setSearchKeyword] = useState('');

  const [searchPopupBoxVisible, setSearchPopupBoxVisible] = useState(false);

  const [form] = Form.useForm();
  const [confirmLoading, setConfirmLoading] = useState(false);

  const showSignInModal = () => {
    SetSignInModal(true);
  };

  const showSignUpModal = () => {
    SetSignUpModal(true);
  };

  const { setLang, lang, theme, setTheme } = useContext(GlobalContext);


  useEffect(() => {
    if (router.query.q) {
      setSearchKeyword(router.query.q as string);
    }
  }, []);

  useEffect(() => {
    const handleRouteChange = (url) => {
      const path = url.split('?')[0];
      const query = new URLSearchParams(url.split('?')[1]).get('q');
      if (path === '/search' && query && query !== searchKeyword) {
        setSearchKeyword(query);
      }
    };

    router.events.on('routeChangeComplete', handleRouteChange);
    return () => {
      router.events.off('routeChangeComplete', handleRouteChange);
    };
  }, [searchKeyword]);

  if (!show) {
    return (
      <div className={styles['fixed-settings']}>
        <Settings
          trigger={
            <Button icon={<IconSettings />} type="primary" size="large" />
          }
        />
      </div>
    );
  }

  const handleLogout = () => {
    baxios
      .post('/user/logout')
      .then((response) => {
        setToken(null);
        Message.info(t['navbar.menu.logout.message']);
        dispatch({
          type: 'update-userInfo',
          payload: { ...userInfo, userLoading: false, isLogin: false },
        });
        if (baxios.defaults.headers.common['Authorization']) {
          delete baxios.defaults.headers.common['Authorization'];
        }
        // localStorage.removeItem('userInfo');
        // window.location.pathname = '/';
      })
      .catch((error) => {
        console.error(error);
      });
  };

  const dropList = (
    // <Menu onClickMenuItem={onMenuItemClick}>
    <Menu>
      {/*<Menu.SubMenu*/}
      {/*  key="role"*/}
      {/*  title={*/}
      {/*    <>*/}
      {/*      <IconUser className={styles['dropdown-icon']} />*/}
      {/*      <span className={styles['user-role']}>*/}
      {/*        {role === 'admin'*/}
      {/*          ? t['menu.user.role.admin']*/}
      {/*          : t['menu.user.role.user']}*/}
      {/*      </span>*/}
      {/*    </>*/}
      {/*  }*/}
      {/*>*/}
      {/*  <Menu.Item onClick={handleChangeRole} key="switch role">*/}
      {/*    <IconTag className={styles['dropdown-icon']} />*/}
      {/*    {t['menu.user.switchRoles']}*/}
      {/*  </Menu.Item>*/}
      {/*</Menu.SubMenu>*/}
      {/*<Menu.Item key="setting">*/}
      {/*  <IconSettings className={styles['dropdown-icon']} />*/}
      {/*  {t['menu.user.setting']}*/}
      {/*</Menu.Item>*/}
      <Menu.Item key="user" onClick={() => {
        router.push({
          pathname: `/user/self`,
        });
      }}>
        <IconUser className={styles['dropdown-icon']} />
        {t['menu.user']}
      </Menu.Item>
      <Menu.Item key="deposit" onClick={() => {
        router.push({
          pathname: `/deposit`,
        });
      }}>
        <IconPaymentCheck className={styles['dropdown-icon']} />
        {t['menu.deposit']}
      </Menu.Item>
      {/*<Menu.SubMenu*/}
      {/*  key="more"*/}
      {/*  title={*/}
      {/*    <div style={{ width: 80 }}>*/}
      {/*      <IconExperiment className={styles['dropdown-icon']} />*/}
      {/*      {t['message.seeMore']}*/}
      {/*    </div>*/}
      {/*  }*/}
      {/*>*/}
      {/*  <Menu.Item key="workplace">*/}
      {/*    <IconDashboard className={styles['dropdown-icon']} />*/}
      {/*    {t['menu.dashboard.workplace']}*/}
      {/*  </Menu.Item>*/}
      {/*  <Menu.Item key="card list">*/}
      {/*    <IconInteraction className={styles['dropdown-icon']} />*/}
      {/*    {t['menu.list.cardList']}*/}
      {/*  </Menu.Item>*/}
      {/*</Menu.SubMenu>*/}

      <Divider style={{ margin: '4px 0' }} />
      <Menu.Item key="logout" onClick={handleLogout}>
        <IconPoweroff className={styles['dropdown-icon']} />
        {t['navbar.logout']}
      </Menu.Item>
    </Menu>
  );

  function onSignInOk() {
    form
      .validate(['username', 'password'])
      .then((res) => {
        console.log(res);
        // const params = {
        //   username: res.username,
        //   password: res.password
        //
        setConfirmLoading(true);
        // sleep 1000ms
        setTimeout(() => {
          const params = new FormData();
          params.append('username', res.username);
          params.append('password', res.password);
          baxios
            .post('/user/login', params)
            .then((response) => {
              const data = response.data;
              if (data.status !== 200) {
                console.error(data.err_msg);
                Message.error(data.err_msg);
                return;
              }
              console.log(data);
              setToken(data.data.token);
              baxios.defaults.headers.common['Authorization'] = `Bearer ${data.data.token}`;
              // localStorage.setItem('userInfo', JSON.stringify(data));
              // window.location.pathname = '/';
              Message.success(t['navbar.model.signin.message.success']);
              SetSignInModal(false);
              // get user info
              FetchUserInfo(dispatch);

            })
            .catch((error) => {
              console.error(error);
              Message.error(error);
            })
            .finally(() => {setConfirmLoading(false);
            });
        }, 1000);
      })
      .catch();
  }
  function onSignUpOk() {
    form
      .validate()
      .then((res) => {
        console.log(res);
        // const params = {
        //   username: res.username,
        //   password: res.password,
        //   nickname: res.nickname,
        // };
        setConfirmLoading(true);
        // sleep 1000ms
        setTimeout(() => {
          const params = new FormData();
          params.append('username', res.username);
          params.append('nickname', res.nickname);
          params.append('password', res.password);
          baxios
            .post('/user/signup', params)
            .then((response) => {
              const data = response.data;
              if (data.status !== 200) {
                console.error(data.err_msg);
                Message.error(data.err_msg);
                return;
              }
              console.log(data);
              setToken(data.data.token);
              baxios.defaults.headers.common['Authorization'] = `Bearer ${data.data.token}`;
              // localStorage.setItem('userInfo', JSON.stringify(data));
              // window.location.pathname = '/';
              Message.success(t['navbar.model.signup.message.success']);
              SetSignUpModal(false);
              SetSignInModal(true);
            })
            .catch((error) => {
              console.error(error);
              Message.error(error);
            })
            .finally(() => {setConfirmLoading(false)});
        }, 1000);
      })
      .catch();
  }

  const handleSearchSubmit = () => {
    // console.log(searchKeyword)
    router.push({
      pathname: '/search',
      query: {
        q: searchKeyword,
      },
    });
    // close popup box
    setSearchPopupBoxVisible(false);
  };

  const formItemLayout = {
    labelCol: {
      span: 4,
    },
    wrapperCol: {
      span: 20,
    },
  };

  return (
    <div className={styles.navbar}>
      <div className={styles.left}>
        <div className={styles.logo}>
          <Link href={"/"}>
            <>
              <Logo />
            </>
          </Link>
        </div>
      </div>
      <ul className={styles.right}>
        <li className={styles['search']}>
          <SearchPopupBox
            searchPopupBoxVisible={searchPopupBoxVisible}
            setSearchPopupBoxVisible={setSearchPopupBoxVisible}
          >
            <Input.Search
              // after enter, popupBox should be closed
              value={searchKeyword}
              onChange={(e) => setSearchKeyword(e)}
              className={styles.round}
              placeholder={t['navbar.search.placeholder']}
              onSearch={searchKeyword !== '' && handleSearchSubmit}
            />
          </SearchPopupBox>
        </li>
        <li>
          {/* upload video button */}
          <Button
            icon={<IconUpload />}
            type="primary"
            shape="round"
            onClick={() => {
              router.push({
                pathname: '/upload',
              });
            }}
          >
            上传短视频
          </Button>
        </li>
        <li>
          <Select
            triggerElement={<IconButton icon={<IconLanguage />} />}
            options={[
              { label: '中文', value: 'zh-CN' },
              { label: 'English', value: 'en-US' },
            ]}
            value={lang}
            triggerProps={{
              autoAlignPopupWidth: false,
              autoAlignPopupMinWidth: true,
              position: 'br',
            }}
            trigger="hover"
            onChange={(value) => {
              setLang(value);
              const nextLang = defaultLocale[value];
              Message.info(`${nextLang['message.lang.tips']}${value}`);
            }}
          />
        </li>
        {/*<li>*/}
        {/*  <MessageBox>*/}
        {/*    <IconButton icon={<IconNotification />} />*/}
        {/*  </MessageBox>*/}
        {/*</li>*/}
        {/* <li>
          <Tooltip
            content={
              theme === 'light'
                ? t['settings.navbar.theme.toDark']
                : t['settings.navbar.theme.toLight']
            }
          >
            <IconButton
              icon={theme !== 'dark' ? <IconMoonFill /> : <IconSunFill />}
              onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}
            />
          </Tooltip>
        </li> */}
        {/* <Settings /> */}
        {isLogin ? (
          <li>
            <Dropdown droplist={dropList} position="br" disabled={userLoading}>
              <Avatar
                size={32}
                style={{ cursor: 'pointer' }}
                onClick={() => {
                  router.push({
                    pathname: `/user/self`,
                  });
                }}
              >
                {userLoading ? (
                  <IconLoading />
                ) : (
                  <Avatar size={40}>
                    {userInfo.avatar_url ? (
                      <img src={userInfo.avatar_url}  alt={null}/>
                    ) : (
                      userInfo.nickname
                    )}
                  </Avatar>
                  // <img alt="avatar" src={userInfo.avatar} />
                )}
              </Avatar>
            </Dropdown>
          </li>
        ) : (
          <li>
            <Button type="text" onClick={showSignInModal}>
              {t['navbar.button.signin']}
            </Button>
            <Button type="text" onClick={showSignUpModal}>
              {t['navbar.button.signup']}
            </Button>
          </li>
        )}
      </ul>
      <Modal
        title={
          <div>
            <h1 className="h1 text-white">{t['navbar.model.signin.title']}</h1>
          </div>
        }
        className={styles['auth-model']}
        visible={signInModal}
        confirmLoading={confirmLoading}
        okText={t['navbar.button.signin']}
        onOk={() => onSignInOk()}
        onCancel={() => SetSignInModal(false)}
        autoFocus={false}
        focusLock={true}
        closable={false}
        simple={true}
      >
        <div>
          <Form
            {...formItemLayout}
            form={form}
            labelCol={{
              style: { flexBasis: 90 },
            }}
            wrapperCol={{
              // style: { flexBasis: 'calc(100% - 90px)' },
            }}
          >
            <FormItem field="username" rules={[{ required: true, message: t['navbar.model.username.require'] }]}>
              <Input
                autoComplete={'off'}
                prefix={<IconUser />}
                placeholder={t['navbar.model.username']}
              />
            </FormItem>
            <FormItem field="password" rules={[{ required: true, message: t['navbar.model.password.require'] }]}>
              <Input.Password
                prefix={<IconLock />}
                placeholder={t['navbar.model.password']}
              />
            </FormItem>
          </Form>
        </div>
      </Modal>
      <Modal
        title={
          <div>
            <h1 className="h1 text-white">{t['navbar.model.signup.title']}</h1>
          </div>
        }
        className={styles['auth-model']}
        visible={signUpModal}
        confirmLoading={confirmLoading}
        okText={t['navbar.button.signup']}
        onOk={() => onSignUpOk()}
        onCancel={() => SetSignUpModal(false)}
        autoFocus={false}
        focusLock={true}
        closable={false}
        simple={true}
      >
        <div>
          <Form
            {...formItemLayout}
            form={form}
            labelCol={{
              style: { flexBasis: 90 },
            }}
            wrapperCol={{
              // style: { flexBasis: 'calc(100% - 90px)' },
            }}
          >
            <FormItem field="username" rules={[{ required: true, message: t['navbar.model.username.require'] }]}>
              <Input
                autoComplete={'off'}
                prefix={<IconUser />}
                placeholder={t['navbar.model.username']}
              />
            </FormItem>
            <FormItem field="nickname" rules={[{ required: true, message: t['navbar.model.nickname.require'] }]}>
              <Input
                autoComplete={'off'}
                prefix={<IconPen />}
                placeholder={t['navbar.model.nickname']}
              />
            </FormItem>
            <FormItem field="password" rules={[{ required: true, message: t['navbar.model.password.require'] }]}>
              <Input.Password
                prefix={<IconLock />}
                placeholder={t['navbar.model.password']}
              />
            </FormItem>

            <FormItem
              field="confirm_password"
              required={true}
              // dependencies={['password']}
              rules={[
                {
                  validator: (v, cb) => {
                    if (!v) {
                      return cb(t['navbar.model.password.confirm.require']);
                    } else if (form.getFieldValue('password') !== v) {
                      return cb(t['navbar.model.password.confirm.error']);
                    }
                    cb(null);
                  },
                },
              ]}
            >
              <Input.Password
                prefix={<IconLock />}
                placeholder={t['navbar.model.password.confirm']}
              />
            </FormItem>
          </Form>
        </div>
      </Modal>
    </div>
  );
}

export default Navbar;
