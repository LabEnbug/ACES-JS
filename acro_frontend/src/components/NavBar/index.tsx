import React, { useContext, useEffect, useState } from 'react';
import {
  Tooltip,
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
  Trigger,
  Skeleton,
} from '@arco-design/web-react';
import {
  IconLanguage,
  IconNotification,
  IconSunFill,
  IconMoonFill,
  IconUser,
  IconSettings,
  IconPoweroff,
  IconExperiment,
  IconDashboard,
  IconInteraction,
  IconTag,
  IconLoading,
  IconLock,
  IconShake,
  IconPlusCircle,
  IconPlus,
  IconUpload,
} from '@arco-design/web-react/icon';
import { useSelector, useDispatch } from 'react-redux';
import store, { GlobalState } from '@/store';
import { GlobalContext } from '@/context';
import useLocale from '@/utils/useLocale';
import Logo from '@/assets/logo.svg';
import MessageBox from '@/components/MessageBox';
import SearchPopupBox from '@/components/SearchPopupBox';
import IconButton from './IconButton';
import Settings from '../Settings';
import styles from './style/index.module.less';
import defaultLocale from '@/locale';
import useStorage from '@/utils/useStorage';
import { generatePermission } from '@/routes';
import { useRouter } from 'next/router';
import GetAxios from '@/utils/getaxios';
import cs from 'classnames';

const FormItem = Form.Item;

function Navbar({ show }: { show: boolean }) {
  const t = useLocale();
  const { userInfo, userLoading } = useSelector((state: GlobalState) => state);
  const dispatch = useDispatch();
  const [_, setUserStatus] = useStorage('userStatus');
  const [role, setRole] = useStorage('userRole', 'admin');
  const [signinmodel, SetSignInModel] = useState(false);
  const [signupmodel, SetSignUpModel] = useState(false);

  const router = useRouter();
  const [searchKeyword, setSearchKeyword] = useState('');

  const [searchPopupBoxVisible, setSearchPopupBoxVisible] = useState(false);

  const showSignInModal = () => {
    SetSignInModel(true);
  };

  const showSignUpModal = () => {
    SetSignUpModel(true);
  };

  const { setLang, lang, theme, setTheme } = useContext(GlobalContext);

  // function logout() {
  //   setUserStatus('logout');
  //   window.location.href = '/login';
  // }

  // function onMenuItemClick(key) {
  //   if (key === 'logout') {
  //     logout();
  //   } else {
  //     Message.info(`You clicked ${key}`);
  //   }
  // }

  useEffect(() => {
    dispatch({
      type: 'update-userInfo',
      payload: {
        userInfo: {
          ...userInfo,
          permissions: generatePermission(role),
        },
      },
    });
  }, [role]);

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

  const handleChangeRole = () => {
    const newRole = role === 'admin' ? 'user' : 'admin';
    setRole(newRole);
  };

  const handlelogout = () => {
    const baxios = GetAxios();

    baxios
      .get('/v1-api/v1/user/logout')
      .then((response) => {
        Message.info(t['navbar.menu.logout.message']);
        localStorage.removeItem('userInfo');
        window.location.pathname = '/';
      })
      .catch((error) => {
        console.error(error);
      });
  };

  const droplist = (
    // <Menu onClickMenuItem={onMenuItemClick}>
    <Menu>
      <Menu.SubMenu
        key="role"
        title={
          <>
            <IconUser className={styles['dropdown-icon']} />
            <span className={styles['user-role']}>
              {role === 'admin'
                ? t['menu.user.role.admin']
                : t['menu.user.role.user']}
            </span>
          </>
        }
      >
        <Menu.Item onClick={handleChangeRole} key="switch role">
          <IconTag className={styles['dropdown-icon']} />
          {t['menu.user.switchRoles']}
        </Menu.Item>
      </Menu.SubMenu>
      <Menu.Item key="setting">
        <IconSettings className={styles['dropdown-icon']} />
        {t['menu.user.setting']}
      </Menu.Item>
      <Menu.SubMenu
        key="more"
        title={
          <div style={{ width: 80 }}>
            <IconExperiment className={styles['dropdown-icon']} />
            {t['message.seeMore']}
          </div>
        }
      >
        <Menu.Item key="workplace">
          <IconDashboard className={styles['dropdown-icon']} />
          {t['menu.dashboard.workplace']}
        </Menu.Item>
        <Menu.Item key="card list">
          <IconInteraction className={styles['dropdown-icon']} />
          {t['menu.list.cardList']}
        </Menu.Item>
      </Menu.SubMenu>

      <Divider style={{ margin: '4px 0' }} />
      <Menu.Item key="logout" onClick={handlelogout}>
        <IconPoweroff className={styles['dropdown-icon']} />
        {t['navbar.logout']}
      </Menu.Item>
    </Menu>
  );

  const [form] = Form.useForm();
  const [confirmLoading, setConfirmLoading] = useState(false);

  function onSignInOk() {
    const baxios = GetAxios();
    form
      .validate()
      .then((res) => {
        console.log(res);
        // const params = {
        //   username: res.username,
        //   password: res.password
        // }
        const params = new FormData();
        params.append('username', res.username);
        params.append('password', res.password);
        baxios
          .post('/v1-api/v1/user/login', params)
          .then((response) => {
            const data = response.data;
            localStorage.setItem('userInfo', JSON.stringify(data));
            window.location.pathname = '/';
          })
          .catch((error) => {
            console.error(error);
          });
      })
      .catch((e) => {
        Message.error(t['navbar.model.signup.message.error']);
      });
  }
  function onSignUpOk() {
    const baxios = GetAxios();
    form
      .validate()
      .then((res) => {
        console.log(res);
        const params = {
          username: res.username,
          password: res.password,
          nickname: res.nickname,
        };
        baxios
          .get('/v1-api/v1/user/signup', { params })
          .then((response) => {
            const data = response.data;
            console.log(data);
            localStorage.setItem('userInfo', JSON.stringify(data));
            window.location.pathname = '/';
          })
          .catch((error) => {
            console.log(1234);
            console.error(error);
          });
      })
      .catch((e) => {
        Message.error(t['navbar.model.signup.message.error']);
      });
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
          <Logo />
          {/*<div className={styles['logo-name']}>ACES 短视频</div>*/}
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
        {userInfo ? (
          <li>
            <Dropdown droplist={droplist} position="br" disabled={userLoading}>
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
                      <img src={userInfo.avatar_url} />
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
        visible={signinmodel}
        onOk={() => onSignInOk()}
        onCancel={() => SetSignInModel(false)}
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
              style: { flexBasis: 'calc(100% - 90px)' },
            }}
          >
            <FormItem field="username" rules={[{ required: true }]}>
              <Input
                prefix={<IconUser />}
                placeholder={t['navbar.model.signin.account']}
              />
            </FormItem>
            <FormItem field="password" rules={[{ required: true }]}>
              <Input.Password
                prefix={<IconLock />}
                placeholder={t['navbar.model.signin.passward']}
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
        visible={signupmodel}
        onOk={() => onSignUpOk()}
        onCancel={() => SetSignUpModel(false)}
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
              style: { flexBasis: 'calc(100% - 90px)' },
            }}
          >
            <FormItem field="username" rules={[{ required: true }]}>
              <Input
                prefix={<IconUser />}
                placeholder={t['navbar.model.signin.account']}
              />
            </FormItem>
            <FormItem field="nickname" rules={[{ required: true }]}>
              <Input
                prefix={<IconUser />}
                placeholder={t['navbar.model.signup.nickname']}
              />
            </FormItem>
            <FormItem field="password" rules={[{ required: true }]}>
              <Input.Password
                prefix={<IconLock />}
                placeholder={t['navbar.model.signin.passward']}
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
                      return cb('confirm_password is required');
                    } else if (form.getFieldValue('password') !== v) {
                      return cb('confirm_password must be equal with password');
                    }
                    cb(null);
                  },
                },
              ]}
            >
              <Input.Password
                prefix={<IconLock />}
                placeholder="please confirm your password"
              />
            </FormItem>
          </Form>
        </div>
      </Modal>
    </div>
  );
}

export default Navbar;
