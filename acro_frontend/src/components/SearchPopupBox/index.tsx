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
import {makeNewSearch} from "@/utils/keywordUtils";
import baxios from "@/utils/getaxios";

function DropContent({setSearchPopupBoxVisible}) {
  const t = useLocale();
  const [loading, setLoading] = useState(false);
  const [hotkeysData, setHotkeysData] = useState([]);
  const [searchHistory, setSearchHistory] = useState([]);

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
    baxios.get('/search/video/hotkeys?' + 'max_count=' + '10')
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
    setSearchHistory(getSearchHistory());
    getHotkeys();
  }, []);


  const showTag = () => {
    return searchHistory.map((item, index) => (
      <Tag key={index} closable color={'gray'}
           style={{ margin: '8px 16px 0 0 ', cursor: "pointer" }}
           onClick={() => {
             router.push({
               pathname: '/search',
               query: {
                 q: item,
               },
             });
             // close popup box
             setSearchPopupBoxVisible(false);
           }}
           onClose={(e) => {
             const searchHistory = getSearchHistory();
             const index = searchHistory.indexOf(item);
             console.log(searchHistory)
             console.log(index)
             if (index > -1) {
               searchHistory.splice(index, 1);
             }
             localStorage.setItem('searchHistory', JSON.stringify(searchHistory));
             e.stopPropagation();
           }}>
        {item}
      </Tag>
    ));
  };

  return (
    <div className={styles['message-box']}>
      <Spin loading={loading} style={{display: 'block'}}>
        <div className={styles['search-popup-column']}>
          <div className={styles['search-popup-title']}>搜索历史</div>
          <div>
            { /* getSearchHistory to tag */}
            {showTag()}
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
                     // close popup box
                     setSearchPopupBoxVisible(false);
                   }}>{item}</div>
              ))}
          </div>
        </div>
      </Spin>
    </div>
  );
}

function SearchPopupBox({children, searchPopupBoxVisible, setSearchPopupBoxVisible}) {
  return (
    <Trigger
      // trigger={'click'}
      onClick={() => setSearchPopupBoxVisible(!searchPopupBoxVisible)}
      onClickOutside={() => setSearchPopupBoxVisible(false)}
      popup={() => <DropContent setSearchPopupBoxVisible={setSearchPopupBoxVisible}/>}
      position="bottom"
      unmountOnExit={true}
      popupAlign={{bottom: 4}}
      popupVisible={searchPopupBoxVisible}
    >
      {/*<Badge count={9} dot>*/}
        {children}
      {/*</Badge>*/}
    </Trigger>
  );
}

export default SearchPopupBox;
