import React, { useState, ReactNode, useRef, useEffect } from 'react';
import { Layout, Menu, Breadcrumb, Spin } from '@arco-design/web-react';
import cs from 'classnames';
import {
  IconUser,
  IconMenuFold,
  IconMenuUnfold,
  IconMusic,
  IconHome,
} from '@arco-design/web-react/icon';
import IconGame from '@/assets/game.svg';
import IconKnowledge from '@/assets/knowledge.svg';
import IconHot from '@/assets/hot.svg';
import IconEntertainment from '@/assets/entertainment.svg';
import IconFantasy from '@/assets/fantasy.svg';
import IconFood from '@/assets/food.svg';
import IconSport from '@/assets/sport.svg';
import IconFashion from '@/assets/fashion.svg';
import { useSelector } from 'react-redux';
import { useRouter } from 'next/router';
import Link from 'next/link';
import qs from 'query-string';
import Navbar from '../components/NavBar';
import Footer from '../components/Footer';
import useRoute, { IRoute } from '@/routes';
import useLocale from '@/utils/useLocale';
import { GlobalState } from '@/store';
import getUrlParams from '@/utils/getUrlParams';
import styles from '@/style/layout.module.less';
import NoAccess from '@/pages/exception/403';
import path from 'path';

const MenuItem = Menu.Item;
const SubMenu = Menu.SubMenu;

const Sider = Layout.Sider;
const Content = Layout.Content;

function getIconFromKey(key) {
  switch (key) {
    case 'video':
      return <IconHome className={styles.icon} />;
    case 'video?type=knowledge':
      return <IconKnowledge className={styles.icon} />;
    case 'video?type=hotpot':
      return <IconHot className={styles.icon} />;
    case 'video?type=game':
      return <IconGame className={styles.icon} />;
    case 'video?type=entertainment':
      return <IconEntertainment className={styles.icon} />;
    case 'video?type=fantasy':
      return <IconFantasy className={styles.icon} />;
    case 'video?type=music':
      return <IconMusic className={styles.icon} />;
    case 'video?type=food':
      return <IconFood className={styles.icon} />;
    case 'video?type=sport':
      return <IconSport className={styles.icon} />;
    case 'video?type=fashion':
      return <IconFashion className={styles.icon} />;
    case 'user':
      return <IconUser className={styles.icon} />;
    default:
      return <div className={styles['icon-empty']} />;
  }
}

function PageLayout({ children }: { children: ReactNode }) {
  const urlParams = getUrlParams();
  const router = useRouter();
  const pathname = router.pathname;
  const query = router.query;
  const currentComponent = qs.parseUrl(pathname).url.slice(1);
  const locale = useLocale();
  const { userInfo, settings, userLoading } = useSelector(
    (state: GlobalState) => state
  );
  const [collapsed, setCollapsed] = useState<boolean>(false);

  const [routes, defaultRoute] = useRoute(userInfo?.permissions);

  const defaultSelectedKeys = [currentComponent || defaultRoute];
  const paths = (currentComponent || defaultRoute).split('/');
  const defaultOpenKeys = paths.slice(0, paths.length - 1);

  const [selectedKeys, setSelectedKeys] =
    useState<string[]>(defaultSelectedKeys);
  const [openKeys, setOpenKeys] = useState<string[]>(defaultOpenKeys);

  const navbarHeight = 60;
  const menuWidth = collapsed ? 48 : settings?.menuWidth;

  const showNavbar = settings?.navbar && urlParams.navbar !== false;
  const showMenu = settings?.menu && urlParams.menu !== false;

  const routeMap = useRef<Map<string, ReactNode[]>>(new Map());
  const menuMap = useRef<
    Map<string, { menuItem?: boolean; subMenu?: boolean }>
  >(new Map());

  const [breadcrumb, setBreadCrumb] = useState([]);

  function onClickMenuItem(key) {
    setSelectedKeys([key]);
  }

  function toggleCollapse() {
    setCollapsed((collapsed) => !collapsed);
  }

  const paddingLeft = showMenu ? { paddingLeft: menuWidth } : {};
  const paddingTop = showNavbar ? { paddingTop: navbarHeight } : {};
  const paddingStyle = { ...paddingLeft, ...paddingTop };

  function renderRoutes(locale) {
    routeMap.current.clear();
    return function travel(_routes: IRoute[], level, parentNode = []) {
      return _routes.map((route) => {
        const { breadcrumb = true, ignore } = route;
        const iconDom = getIconFromKey(route.key);
        const titleDom = (
          <>
            {iconDom} {locale[route.name] || route.name}
          </>
        );

        routeMap.current.set(
          `/${route.key}`,
          breadcrumb ? [...parentNode, route.name] : []
        );

        const visibleChildren = (route.children || []).filter((child) => {
          const { ignore, breadcrumb = true } = child;
          if (ignore || route.ignore) {
            routeMap.current.set(
              `/${child.key}`,
              breadcrumb ? [...parentNode, route.name, child.name] : []
            );
          }

          return !ignore;
        });

        if (ignore) {
          return '';
        }
        if (visibleChildren.length) {
          menuMap.current.set(route.key, { subMenu: true });
          return (
            <SubMenu key={route.key} title={titleDom}>
              {travel(visibleChildren, level + 1, [...parentNode, route.name])}
            </SubMenu>
          );
        }
        menuMap.current.set(route.key, { menuItem: true });
        return (
          <MenuItem key={route.key}>
            <Link href={`/${route.key}`}>
              <a style={{ textDecoration: 'none' }}>{titleDom}</a>
            </Link>
          </MenuItem>
        );
      });
    };
  }

  function updateMenuStatus() {
    const pathKeys = pathname.split('/');
    const newSelectedKeys: string[] = [];
    const newOpenKeys: string[] = [...openKeys];
    while (pathKeys.length > 0) {
      const currentRouteKey = pathKeys.join('/');
      const menuKey = currentRouteKey.replace(/^\//, '');
      const menuType = menuMap.current.get(menuKey);
      if (menuType && menuType.menuItem) {
        newSelectedKeys.push(menuKey);
      }
      if (menuType && menuType.subMenu && !openKeys.includes(menuKey)) {
        newOpenKeys.push(menuKey);
      }
      pathKeys.pop();
    }
    if (pathname.includes('video') && query['type']) {
      setSelectedKeys([`video?type=${query['type']}`]);
    } else setSelectedKeys(newSelectedKeys);
    setOpenKeys(newOpenKeys);
  }

  useEffect(() => {
    const routeConfig = routeMap.current.get(pathname);
    // setBreadCrumb(routeConfig || []);
    updateMenuStatus();
  }, [pathname]);

  return (
    <Layout className={styles.layout}>
      <div
        className={cs(styles['layout-navbar'], {
          [styles['layout-navbar-hidden']]: !showNavbar,
        })}
      >
        <Navbar show={showNavbar} />
      </div>
      {userLoading ? (
        <Spin className={styles['spin']} />
      ) : (
        <Layout>
          {showMenu && (
            <Sider
              className={styles['layout-sider']}
              width={menuWidth}
              collapsed={collapsed}
              onCollapse={setCollapsed}
              trigger={null}
              collapsible
              breakpoint="xl"
              style={paddingTop}
            >
              <div className={styles['menu-wrapper']}>
                <Menu
                  collapse={collapsed}
                  onClickMenuItem={onClickMenuItem}
                  selectedKeys={selectedKeys}
                  openKeys={openKeys}
                  onClickSubMenu={(_, openKeys) => {
                    setOpenKeys(openKeys);
                  }}
                >
                  {renderRoutes(locale)(routes, 1)}
                </Menu>
              </div>
              <div className={styles['info-up-of-collapse-btn']}>
                <div>2023</div>
                <div>@</div>
                <div>ACES</div>
              </div>
              <div className={styles['collapse-btn']} onClick={toggleCollapse}>
                {collapsed ? <IconMenuUnfold /> : <IconMenuFold />}
              </div>
            </Sider>
          )}
          <Layout className={styles['layout-content']} style={paddingStyle}>
            <div
              className={
                pathname !== '/video' ? styles['layout-content-wrapper'] : null
              }
            >
              {!!breadcrumb.length && (
                <div className={styles['layout-breadcrumb']}>
                  <Breadcrumb>
                    {breadcrumb.map((node, index) => (
                      <Breadcrumb.Item key={index}>
                        {typeof node === 'string' ? locale[node] || node : node}
                      </Breadcrumb.Item>
                    ))}
                  </Breadcrumb>
                </div>
              )}
              <Content>
                {routeMap.current.has(pathname) ? children : <NoAccess />}
              </Content>
            </div>
          </Layout>
        </Layout>
      )}
    </Layout>
  );
}

export default PageLayout;
