import React from "react";
import {  Tabs, Typography, Comment, Avatar, Input, Tooltip, Message, Button } from '@arco-design/web-react';
import CommentTab from './Comment';
import RelatedVideosTab from './RelatedVideos';
import locale from './locale';
import useLocale from '@/utils/useLocale';

const TextArea = Input.TextArea;
const TabPane = Tabs.TabPane;


function SideBar(props) {
    const {videoinfo} = props;
    const t = useLocale(locale);
    return (
        <Tabs defaultActiveTab='1'>
            <TabPane key='1' title={t['comment']} style={{'color': '#ffffff'}} >
                <CommentTab videoinfo={videoinfo} />
            </TabPane>
            <TabPane key='2' title={t['related.video']} style={{'color': '#ffffff'}} >
              <RelatedVideosTab videoInfo={videoinfo} />
            </TabPane>
        </Tabs>
      );
  }
  
  export default SideBar;