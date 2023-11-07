import React, {useEffect, useRef, useState} from "react";
import {Button, Empty, List} from "@arco-design/web-react";
import styles from './style/index.module.less';
import CardBlock from "./card-block"
import useLocale from "@/utils/useLocale";
import locale from "./locale";
import baxios from "@/utils/getaxios";
import {useRouter} from "next/router";

const defaultVideoList = new Array(0).fill({});
function RelatedVideos (props) {
  const { videoInfo } = props;
  const t = useLocale(locale);
  const tg = useLocale();
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const listRef = useRef(null);
  const [relatedVideoData, setRelatedVideoData] = useState(defaultVideoList);
  const [isEndData, setIsEndData] = useState(false);

  const getVideoData = () => {
    setIsEndData(false);
    setLoading(true);
    // sleep
    // await new Promise(resolve => setTimeout(resolve, 3000));
    baxios
      .get('/videos/' + videoInfo['video_uid'] + '/related' + '?' +
        'limit=12')
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data['err_msg']);
          return;
        }
        setRelatedVideoData(data.data['video_list']);
      })
      .catch((error) => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  };

  const getMoreData = () => {
    setLoading(true);
    baxios
      .get('/videos/' + videoInfo['video_uid'] + '/related' + '?' +
        'start=' + relatedVideoData.length.toString() + '&' +
        'limit=12')
      .then((response) => {
        const data = response.data;
        if (data.status !== 200) {
          console.error(data.err_msg);
          setIsEndData(true);
          return;
        }
        setRelatedVideoData(relatedVideoData.concat(data.data.video_list));
      })
      .catch((error) => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    if (router.isReady) {
      setRelatedVideoData(defaultVideoList);
      getVideoData();
    }
  }, [router.isReady, videoInfo.video_uid]);


  const ContentContainer = ({ children }) => (
    <div style={{ textAlign: 'center', marginTop: 4 }}>
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        {children}
      </div>
    </div>
  );

  const LoadingIndicator = () => (
    <div style={{ textAlign: 'center', marginTop: 14 }}>
      <span style={{ color: 'var(--color-text-3)' }}>
        {/*<IconLoading style={{marginRight: 8, color: 'rgb(var(--arcoblue-6))'}}/>*/}
        {t['card.loading']}
      </span>
    </div>
  );

  return (
    <List
      ref={listRef}
      className={styles['card-list']}
      grid={{
        xs: 12,
        sm: 12,
        md: 12,
        lg: 12,
        xl: 12,
        xxl: 12,
      }}
      noDataElement={
        loading ? (
          <div />
        ) : (
          <Empty
            description={
              <ContentContainer>
                          <span
                            style={{
                              color: 'var(--color-text-3)',
                              marginTop: '16px',
                            }}
                          >
                            {t['related.no.video']}
                          </span>
              </ContentContainer>
            }
          ></Empty>
        )
      }
      style={{ overflowY: 'scroll', height: 'calc(100vh - 135px)' }}
      dataSource={relatedVideoData}
      bordered={false}
      onReachBottom={() => {
        getMoreData();
      }}
      render={(item, _) => (
        <List.Item style={{ padding: '4px 4px' }}>
          {
            <CardBlock
              card={item}
              loading={loading}
            />
          }
        </List.Item>
      )}
      loading={loading}
      offsetBottom={300}
      footer={
        loading ? (
          <LoadingIndicator />
        ) : isEndData ? (
          <ContentContainer>
                      <span
                        style={{
                          color: 'var(--color-text-3)',
                          marginBottom: '4px',
                        }}
                      >
                        {t['card.noMoreContent']}
                      </span>
          </ContentContainer>
        ) : (
          relatedVideoData.length !== 0 && (
            <ContentContainer>
              <Button
                type="text"
                style={{ color: 'var(--color-text-1)' }}
                onClick={() => {
                  getMoreData();
                }}
              >
                {`${t['card.loadMore']}${t['related.video']}`}
              </Button>
            </ContentContainer>
          )
        )
      }
    />
  );
}

export default RelatedVideos;