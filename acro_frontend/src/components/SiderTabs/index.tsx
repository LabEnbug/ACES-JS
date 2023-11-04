import { useEffect, useState, useRef } from 'react';
import {  Tabs, Typography, Comment, Avatar, Input, Tooltip, Message, Button } from '@arco-design/web-react';
import CommentTab from './Comment';
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
        </Tabs>
      );
  }
  
  export default SideBar;