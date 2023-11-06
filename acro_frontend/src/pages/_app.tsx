import React, { useEffect, useMemo } from 'react';
import { useRouter } from 'next/router';
import cookies from 'next-cookies';
import Head from 'next/head';
import type { AppProps } from 'next/app';
import { createStore } from 'redux';
import {Provider, useSelector} from 'react-redux';
import '../style/global.less';
import { ConfigProvider } from '@arco-design/web-react';
import zhCN from '@arco-design/web-react/es/locale/zh-CN';
import enUS from '@arco-design/web-react/es/locale/en-US';
import NProgress from 'nprogress';
import rootReducer from '../store';
import { GlobalContext } from '../context';
import changeTheme from '@/utils/changeTheme';
import useStorage from '@/utils/useStorage';
import Layout from './layout';
import {getToken} from "@/utils/authentication";
import baxios from "@/utils/getaxios";
import FetchUserInfo from "@/utils/getuserinfo";

const store = createStore(rootReducer);

interface RenderConfig {
  arcoLang?: string;
  arcoTheme?: string;
}

export default function MyApp({
  pageProps,
  Component,
  renderConfig,
}: AppProps & { renderConfig: RenderConfig }) {
  const { arcoLang, arcoTheme } = renderConfig;
  const [lang, setLang] = useStorage('arco-lang', arcoLang || 'zh-CN');
  const [theme, setTheme] = useStorage('arco-theme', 'dark');
  // setTheme('dark')
  const router = useRouter();
  const locale = useMemo(() => {
    switch (lang) {
      case 'zh-CN':
        return zhCN;
      case 'en-US':
        return enUS;
      default:
        return zhCN;
    }
  }, [lang]);

  // async function GetUserInfo() {
  //   try {
  //     const response = await baxios
  //       .post('/user/info');
  //     const data = response.data;
  //     if (data.status !== 200) {
  //       console.error(data.err_msg);
  //       throw new Error(data.err_msg);
  //     }
  //     console.log(data);
  //     return data.data.user;
  //   } catch (error) {
  //     throw error;
  //   }
  // }
  //
  //
  // // how to make fetchUserInfo() usable by other components?
  // function fetchUserInfo() {
  //   store.dispatch({
  //     type: 'update-userInfo',
  //     payload: { userLoading: true },
  //   });
  //   GetUserInfo().then((userinfo) => {
  //     store.dispatch({
  //       type: 'update-userInfo',
  //       payload: { userInfo: userinfo, userLoading: false, isLogin: true },
  //     });
  //   }).catch((error) => {
  //     store.dispatch({
  //       type: 'update-userInfo',
  //       payload: { userInfo: null, userLoading: false, isLogin: false },
  //     });
  //   });
  // }

  useEffect(() => {
    const token = getToken();
    if (token) {
      baxios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    }
    store.dispatch({
      type: 'init',
      payload: {init: true},
    });
    FetchUserInfo(store.dispatch);
    // if (checkLogin()) {
    //   fetchUserInfo();
    // }
    // else if (window.location.pathname.replace(/\//g, '') !== 'login') {
    //   window.location.pathname = '/login';
    // }
  }, []);

  useEffect(() => {
    const handleStart = () => {
      NProgress.set(0.4);
      NProgress.start();
    };

    const handleStop = () => {
      NProgress.done();
    };

    router.events.on('routeChangeStart', handleStart);
    router.events.on('routeChangeComplete', handleStop);
    router.events.on('routeChangeError', handleStop);

    return () => {
      router.events.off('routeChangeStart', handleStart);
      router.events.off('routeChangeComplete', handleStop);
      router.events.off('routeChangeError', handleStop);
    };
  }, [router]);

  useEffect(() => {
    document.cookie = `arco-lang=${lang}; path=/`;
    document.cookie = `arco-theme=${theme}; path=/`;
    changeTheme(theme);
  }, [lang, theme]);

  const contextValue = {
    lang,
    setLang,
    theme,
    setTheme,
  };
  return (
    <>
      <Head>
        <link
          rel="shortcut icon"
          type="image/x-icon"
          href="https://unpkg.byted-static.com/latest/byted/arco-config/assets/favicon.ico"
        />
        ACES短视频
      </Head>
      <ConfigProvider
        locale={locale}
        componentConfig={{
          Card: {
            bordered: false,
          },
          List: {
            bordered: false,
          },
          Table: {
            border: false,
          },
        }}
      >
        <Provider store={store}>
          <GlobalContext.Provider value={contextValue}>
            {Component.displayName === 'LoginPage' ? (
              <Component {...pageProps} suppressHydrationWarning />
            ) : (
              <Layout>
                <Component {...pageProps} suppressHydrationWarning />
              </Layout>
            )}
          </GlobalContext.Provider>
        </Provider>
      </ConfigProvider>
    </>
  );
}

// fix: next build ssr can't attach the localstorage
MyApp.getInitialProps = async (appContext) => {
  const { ctx } = appContext;
  const serverCookies = cookies(ctx);
  return {
    renderConfig: {
      arcoLang: serverCookies['arco-lang'],
      arcoTheme: serverCookies['arco-theme'],
    },
  };
};
