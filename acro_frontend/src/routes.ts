import auth, { AuthParams } from '@/utils/authentication';
import { useEffect, useMemo, useState } from 'react';

export type IRoute = AuthParams & {
  name: string;
  key: string;
  // 当前页是否展示面包屑
  breadcrumb?: boolean;
  children?: IRoute[];
  // 当前路由是否渲染菜单项，为 true 的话不会在菜单中显示，但可通过路由地址访问。
  ignore?: boolean;
};

export const routes: IRoute[] = [
  {
    name: 'menu.video.comprehensive',
    key: 'video',
  },
  {
    name: 'menu.video.knowledge',
    key: 'video?type=knowledge',
  },
  {
    name: 'menu.video.hotpot',
    key: 'video?type=hotpot',
  },
  {
    name: 'menu.video.game',
    key: 'video?type=game',
  },
  {
    name: 'menu.video.entertainment',
    key: 'video?type=entertainment',
  },
  {
    name: 'menu.video.fantasy',
    key: 'video?type=fantasy',
  },
  {
    name: 'menu.video.music',
    key: 'video?type=music',
  },
  {
    name: 'menu.video.food',
    key: 'video?type=food',
  },
  {
    name: 'menu.video.sport',
    key: 'video?type=sport',
  },
  {
    name: 'menu.video.fashion',
    key: 'video?type=fashion',
  },
  {
    name: 'menu.search',
    key: 'search',
    ignore: true,
  },
  {
    name: 'menu.user',
    key: 'user/[username]',
    ignore: true,
  },
  {
    name: 'menu.upload',
    key: 'upload',
    ignore: true,
  },
  {
    name: 'menu.edit',
    key: 'edit',
    ignore: true,
  },
  {
    name: 'menu.deposit',
    key: 'deposit',
    ignore: true,
  },
  {
    name: 'menu.promotion',
    key: 'promote',
    ignore: true,
  },
  {
    name: 'menu.advertisement',
    key: 'advertise',
    ignore: true,
  },
  {
    name: 'menu.exception',
    key: 'exception',
    ignore: true,
    children: [
      {
        name: 'menu.exception.403',
        key: 'exception/403',
      },
      {
        name: 'menu.exception.404',
        key: 'exception/404',
      },
      {
        name: 'menu.exception.500',
        key: 'exception/500',
      },
    ],
  },
];

export const getName = (path: string, routes) => {
  return routes.find((item) => {
    const itemPath = `/${item.key}`;
    if (path === itemPath) {
      return item.name;
    } else if (item.children) {
      return getName(path, item.children);
    }
  });
};

export const generatePermission = (role: string) => {
  const actions = role === 'admin' ? ['*'] : ['read'];
  const result = {};
  routes.forEach((item) => {
    if (item.children) {
      item.children.forEach((child) => {
        result[child.name] = actions;
      });
    }
  });
  return result;
};

const useRoute = (userPermission): [IRoute[], string] => {
  const filterRoute = (routes: IRoute[], arr = []): IRoute[] => {
    if (!routes.length) {
      return [];
    }
    for (const route of routes) {
      const { requiredPermissions, oneOfPerm } = route;
      let visible = true;
      if (requiredPermissions) {
        // visible = auth({ requiredPermissions, oneOfPerm }, userPermission);
        visible = true;
      }

      if (!visible) {
        continue;
      }
      if (route.children && route.children.length) {
        const newRoute = { ...route, children: [] };
        filterRoute(route.children, newRoute.children);
        if (newRoute.children.length) {
          arr.push(newRoute);
        }
      } else {
        arr.push({ ...route });
      }
    }

    return arr;
  };

  const [permissionRoute, setPermissionRoute] = useState(routes);

  useEffect(() => {
    const newRoutes = filterRoute(routes);
    setPermissionRoute(newRoutes);
  }, [JSON.stringify(userPermission)]);

  const defaultRoute = useMemo(() => {
    const first = permissionRoute[0];
    if (first) {
      const firstRoute = first?.children?.[0]?.key || first.key;
      return firstRoute;
    }
    return '';
  }, [permissionRoute]);

  return [permissionRoute, defaultRoute];
};

export default useRoute;
