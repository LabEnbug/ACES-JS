import { useEffect, useState, useRef } from 'react';
import {  Tabs, Typography, Comment, Avatar, Input, Tooltip, Message, Button } from '@arco-design/web-react';
import useLocale from '@/utils/useLocale';
import locale from './locale';
import { IconHeart, IconMessage, IconStar } from '@arco-design/web-react/icon';
import styles from './style/index.module.less';
import GetUserInfo from "@/utils/getuserinfo";
import cs from 'classnames';
import { VideoOne } from '@icon-park/react';
import baxios from "@/utils/getaxios";
const TextArea = Input.TextArea;
const TabPane = Tabs.TabPane;



function getNowFormatDate() {
    let date = new Date(),
    year = date.getFullYear(), //获取完整的年份(4位)
    month = date.getMonth() + 1, //获取当前月份(0-11,0代表1月)
    strDate = date.getDate() // 获取当前日(1-31)
    if (month < 10) month = `0${month}` // 如果月份是个位数，在前面补0
    if (strDate < 10) strDate = `0${strDate}` // 如果日是个位数，在前面补0
    return `${year}-${month}-${strDate}T`
}

const JudgeStatus = (data: any) => {
    if (data.status != 200) {
      // Message.error(t['message.notfind'])
      return false;
    }
    return true;
}
const actions = (props)=> {
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
      param.append('video_uid', uid);
      param.append('content', e.target.value);
      param.append('quote_comment_id', comment_id.toString());
      baxios.post('v1-api/v1/video/comment/make', param).then(res=> {
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
            console.log(pre)
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
                <TextArea value={value} className={styles['replay-input']} onChange={(e)=>{ setValue(e)}} onKeyDown={(e)=>{handleKeyDownBottom(e, video_uid, comment_id, addC, setP)}}  placeholder="写下你的回复..." />
             </Tooltip>
        </div>
      )}
    </div>)
}
  
export default actions;