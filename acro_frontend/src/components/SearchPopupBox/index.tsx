import React, {useEffect, useState} from 'react';
import axios from 'axios';
import {
  Trigger,
  Badge,
  Tabs,
  Avatar,
  Spin,
  Button, Tag, Divider,
} from '@arco-design/web-react';
import useLocale from '../../utils/useLocale';
import styles from './style/index.module.less';
import {useRouter} from "next/router";
import {IconRefresh} from "@arco-design/web-react/icon";

function DropContent() {
  const t = useLocale();
  const [loading, setLoading] = useState(false);
  const [hotkeysData, setHotkeysData] = useState([]);

  const router = useRouter();

  // get search history from storage
  const getSearchHistory = () => {
    const searchHistory = localStorage.getItem('searchHistory');
    if (searchHistory) {
      return JSON.parse(searchHistory);
    }
    return [];
  }

  // get hotkeys
  const getHotkeys = () => {
    setLoading(true);
    const params = new FormData();
    params.append('max_count', '10');
    axios.post('/v1-api/v1/video/search/hotkeys', params)
      .then(response => {
        const data = response.data
        if (data.status !== 200) {
          console.error(data.message);
          return;
        }
        console.log(data.data)
        setHotkeysData(data.data.hotkeys);
      })
      .catch(error => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  }

  useEffect(() => {
    getSearchHistory();
    getHotkeys();
  }, []);


  return (
    <div className={styles['message-box']}>
      <Spin loading={loading} style={{display: 'block'}}>
        <div className={styles['search-popup-column']}>
          <div className={styles['search-popup-title']}>搜索历史</div>
          <div>
            { /* getSearchHistory to tag */}
            {getSearchHistory().map((item, index) => (
              <Tag key={index} closable color={'gray'}
                   style={{ margin: '8px 16px 0 0 ', cursor: "pointer" }}
                   onClick={() => {
                     router.push({
                       pathname: '/search',
                       query: {
                         q: item,
                       },
                     });
                   }}
                   onClose={() => {
                      const searchHistory = getSearchHistory();
                      const index = searchHistory.indexOf(item);
                      if (index > -1) {
                        searchHistory.splice(index, 1);
                      }
                      localStorage.setItem('searchHistory', JSON.stringify(searchHistory));
                   }}>
                {item}
              </Tag>
            ))}
          </div>
        </div>
        <Divider style={{ margin: '4px' }}></Divider>
        <div className={styles['search-popup-column']}>
          <div className={styles['search-popup-title']}>
            猜你想搜
            <div className={styles['search-popup-change-hotkeys']} onClick={getHotkeys}><IconRefresh />换一换</div>
          </div>
          <div>
            {hotkeysData.map((item, index) => (
              <div className={styles['hotkey-item']} key={index}
                   onClick={() => {
                     router.push({
                       pathname: '/search',
                       query: {
                         q: item,
                       },
                     });
                   }}>{item}</div>
              ))}
          </div>
        </div>
      </Spin>
    </div>
  );
}

function SearchPopupBox({children}) {
  return (
    <Trigger
      trigger="click"
      popup={() => <DropContent/>}
      position="bottom"
      unmountOnExit={false}
      popupAlign={{bottom: 4}}
    >
      <Badge count={9} dot>
        {children}
      </Badge>
    </Trigger>
  );
}

export default SearchPopupBox;
