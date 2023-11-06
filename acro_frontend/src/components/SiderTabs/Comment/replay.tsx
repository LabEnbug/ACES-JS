import React from "react";
import { useEffect, useState, useRef } from 'react';
import {  Tabs, Typography, Comment, Avatar, Input, Tooltip, Message, Button } from '@arco-design/web-react';
import useLocale from '@/utils/useLocale';
import locale from './locale';
import { IconHeart, IconMessage, IconStar } from '@arco-design/web-react/icon';
import styles from './style/index.module.less';
import baxios from "@/utils/getaxios";
import { useSelector, useDispatch } from 'react-redux';
import store, { GlobalState } from '@/store';


const TextArea = Input.TextArea;
const TabPane = Tabs.TabPane;

const JudgeStatus = (data: any) => {
    if (data.status != 200) {
      // Message.error(t['message.notfind'])
      return false;
    }
    return true;
}
const Actions = (props)=> {
  const { isLogin } = useSelector((state: GlobalState) => state);
  const t = useLocale(locale);
  const {time, comment_id, video_uid, addC, quote_comment, setP} = props;
  const replyInputRef = useRef(null);
  const [showReply, setShowReply] = useState(false);
  const [value, setValue] = useState('');
  const handleReplyClick = () => {
    setShowReply(true);
  };


  const handleKeyDownBottom = (e, uid, comment_id, addC, setP) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      const param = new FormData();
      // 如果只按下了Enter，阻止默认行为并触发你的事件
      e.preventDefault();
      if (e.target.value === '') {
        Message.error(t['comment.input.enter.empty']);
        return;
      }
      param.append('content', e.target.value);
      param.append('quote_comment_id', comment_id.toString());
      baxios.post('/videos/' + uid + '/comments', param).then(res=> {
        if (JudgeStatus(res.data)) {
          const data = res.data.data;
          Message.info(t['comment.input.post.success']);
          addC(pre => {
            return {
                ...pre,
                [data.comment.quote_comment_id]: pre[data.comment.quote_comment_id].concat([{
                    quote_user: quote_comment['child_comment_count_left'] ? null : quote_comment['user'],
                    comment_time: data.comment.comment_time,
                    content: e.target.value,
                    user: data.comment.user,
                    id: data.comment.id
                }])
            }
          });
          setP(pre=>{
            const new_array = [];
            pre.forEach(item=> {
              if (item.id === data.comment.quote_comment_id) {
                item.child_comment_count_left +=1;
              }
              new_array.push(item);
            });
            return new_array;
          });
          setValue('');
          setShowReply(false);
        }
      }).catch(e => {
        Message.info(t['comment.input.post.failed']);
        console.error(e);
      })
    }
  };
  // 点击其他地方隐藏输入框
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (replyInputRef.current && !replyInputRef.current.contains(event.target)) {
        setShowReply(false);
      }
    };

    // 绑定事件监听器
    document.addEventListener('mousedown', handleClickOutside);

    // 组件卸载时移除事件监听器
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  });
  return (
    <div>
      <span onClick={handleReplyClick} style={{cursor:'pointer'}}>
        <IconMessage/> 回复
      </span>
      <span className={styles['comment-time-text']}>{time}</span>
      {showReply && (
        <div ref={replyInputRef}>
            <Tooltip position='tr' trigger='hover' content={t['comment.input.enter']}>
                <TextArea value={value} className={styles['replay-input']} disabled={!isLogin} onChange={(e)=>{ setValue(e)}} onKeyDown={(e)=>{handleKeyDownBottom(e, video_uid, comment_id, addC, setP)}}  placeholder={isLogin ? t['comment.input.placeholder'] :   t['comment.input.placeholder.plslog']} />
             </Tooltip>
        </div>
      )}
    </div>)
}
  
export default Actions;