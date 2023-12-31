import React, { useEffect, useState, useRef } from 'react';
import {  Tabs, Typography, Comment, Avatar, Input, Tooltip, Message, Button } from '@arco-design/web-react';
import useLocale from '@/utils/useLocale';
import locale from './locale';
import styles from './style/index.module.less';
import Replay from './replay';
import baxios from "@/utils/getaxios";
import { useSelector, useDispatch } from 'react-redux';
import store, { GlobalState } from '@/store';
import {parseTime} from "@/utils/timeUtils";
import {useRouter} from "next/router";

const TextArea = Input.TextArea;

function CommentDrawer(props) {
    const {videoinfo} = props;
    const t = useLocale(locale);
    const tg = useLocale();
    const [valuebottom, setValueBottom] = useState('');
    const [displaynomore, setDisplayNoMore] = useState(true);
    const [ comment, SetComment ] = useState([]);
    const [ commentS, SetCommentS ] = useState({});  
    const { isLogin } = useSelector((state: GlobalState) => state);
    const router = useRouter();
    // const [  ]

    const handleScroll = (e) => {
      const atBottom = e.target.scrollHeight - e.target.scrollTop === e.target.clientHeight;
      if (atBottom) {
        // div滚动到底部时的逻辑处理

        baxios.get(
          '/videos/' + videoinfo['video_uid'] + '/comments' + '?' +
          'limit=' + '10' + '&' +
          'start=' + Object.keys(comment).length.toString() + '&' +
          'comment_id=' + '0'
        ).then(res=>{
          if (res.data.status === 200) {
            const news = {}
            res.data.data.comment_list.forEach(item => {
              if (item.child_comment_list && item.child_comment_list.length > 0) {
                news[item.id] = [].concat(item.child_comment_list);
              } else {
                news[item.id] = [];
              } 
            });
            SetCommentS(pre=>(Object.assign({}, news, pre)));
            if (res.data.data.comment_list.length < 20) {
              setDisplayNoMore(true);
            }
            SetComment(pre=>(pre.concat(res.data.data.comment_list)));
          } else {
            setDisplayNoMore(true);
          }   
        }).catch(e=>{
          setDisplayNoMore(true);
        });
      }
    };

    const fetchMoreComment = (comment_id, offset) => {

      baxios.get(
        '/videos/' + videoinfo['video_uid'] + '/comments' + '?' +
        'limit=' + '5' + '&' +
        'start=' + offset.toString() + '&' +
        'comment_id=' + comment_id.toString()
      ).then(res=>{
        if (res.data.status === 200) {
          commentS[comment_id] = commentS[comment_id].concat(res.data.data.child_comment_list);
          
          SetCommentS(pre=> ({
            ...pre,
            comment_id: pre[comment_id].concat(res.data.data.child_comment_list)
          }));
        }
      }).catch(e=>{
        Message.error(t['comment.fetch.failed']);
      });
    }

    const generateSonC = (comment_info, index) => {
      return (
        <Comment
          key={index}
          actions={<Replay time = {parseTime(comment_info['comment_time'], tg)}
                           comment_id={comment_info['id']} 
                           video_uid={videoinfo['video_uid']} 
                           addC={SetCommentS} 
                           quote_comment={comment_info}
                           setP={SetComment}  />}
          author={<> {comment_info['user']['nickname']} { comment_info['quote_user'] ? <><div className={styles['right-arrow']} />  <span> {comment_info['quote_user']['nickname']} </span></> : <></>} </>}
          avatar={(    
            <Avatar
              autoFixFontSize={true}
              style={{
                // backgroundColor: '#000000',
              }}
              onClick={(event) => {
                router.push({
                  pathname: '/user/' + comment_info['user']['username'],
                });
                event.stopPropagation();
              }}
            >
              {comment_info['user']['avatar_url'] ? (
                <img src={comment_info['user']['avatar_url']} alt={null}/>
              ) : (
                comment_info['user']['nickname']
              )}
            </Avatar>)}
          content={<div>{comment_info['content']}</div>}
          // datetime={comment_info['comment_time'].split('T')[0]} 
        />
      )
    }

    const generateParentC = (comment_info, index) => {
      const fetchmore = comment_info['child_comment_count_left'] > commentS[comment_info.id].length && commentS[comment_info.id].length > 0
      return (
          <Comment
                actions={<Replay time = {parseTime(comment_info['comment_time'], tg)}
                                comment_id={comment_info['id']} 
                                video_uid={videoinfo['video_uid']} 
                                addC={SetCommentS} 
                                quote_comment={comment_info}
                                setP={SetComment}  />}
                key={index}
                author={comment_info['user']['nickname']}
                avatar= {(
                  <Avatar
                    autoFixFontSize={true}
                    style={{
                      // backgroundColor: '#000000',
                    }}
                    onClick={(event) => {
                      router.push({
                        pathname: '/user/' + comment_info['user']['username'],
                      });
                      event.stopPropagation();
                    }}
                  >
                    {comment_info['user']['avatar_url'] ? (
                      <img src={comment_info['user']['avatar_url']} alt={null} />
                    ) : (
                      comment_info['user']['nickname']
                    )}
                  </Avatar>)}
                content={<div>{comment_info['content']}</div>}
                // datetime={comment_info['comment_time'].split('T')[0]}
          >
            { 
              commentS[comment_info.id].map((item, index) => (
                // 为每个生成的组件分配一个key，这里使用了item的id作为key
                generateSonC(item, index)
              ))
            }
            {fetchmore&&
            <Button type='text' status='success' onClick={ fetchmore  ? (e)=> {fetchMoreComment(comment_info.id, commentS[comment_info.id].length)} : null}>
              <div className = {styles['comment-div']} /> <span className={styles['comment-div-text']}> {fetchmore ? `展开更多(${comment_info['child_comment_count_left'] - commentS[comment_info.id].length})` : '无更多评论'}</span>
            </Button>
            }
          </Comment>
      )
    }

    const JudgeStatus = (data: any) => {
      if (data.status != 200) {
        Message.error(data.err_msg);
        return false;
      }
      return true;
    }

    const handleKeyDownBottom = (e, uid) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        const param = new FormData();
        // 如果只按下了Enter，阻止默认行为并触发你的事件
        e.preventDefault();
        if (e.target.value === '') {
          Message.error(t['comment.input.enter.empty']);
          return;
        }
        param.append('content', e.target.value);
        param.append('quote_comment_id', '0');
        baxios.post('/videos/' + uid + '/comments', param).then(res=> {
          if (JudgeStatus(res.data)) {
            const data = res.data.data;
            Message.info(t['comment.input.post.success']);
            SetCommentS(pre=>({
              ...pre,
              [data.comment.id]: []
            }));
            SetComment(pre => {
                return [{
                  comment_time: data.comment.comment_time,
                  content: e.target.value,
                  user: data.comment.user,
                  id: data.comment.id,
                  child_comment_count_left: 0,
                }].concat(pre);
            });
            setValueBottom('');
          }
        }).catch(e => {
          Message.info(t['comment.input.post.failed']);
          console.error(e);
        })
      }
    };


    useEffect(()=> {
      SetCommentS({});
      SetComment([]);

      baxios.get(
        '/videos/' + videoinfo['video_uid'] + '/comments' + '?' +
        'limit=' + '20' + '&' +
        'start=' + Object.keys(comment).length.toString() + '&' +
        'comment_id=' + '0'
      ).then(res=>{
        if (res.data.status === 200) {
            const news = {}
            res.data.data.comment_list.forEach(item => {
              if (item.child_comment_list && item.child_comment_list.length > 0) {
                news[item.id] = [].concat(item.child_comment_list);
                item.child_comment_count_left += 1;
              } else {
                news[item.id] = [];
              } 
            });
            SetCommentS(news);
            if (res.data.data.comment_list.length < 20) {
              setDisplayNoMore(true);
            }
            SetComment(pre=>(pre.concat(res.data.data.comment_list)));
        } else {
          setDisplayNoMore(true);
        }   
      }).catch(e=>{
        setDisplayNoMore(true);
      });
    }, [videoinfo.video_uid])

    return (
          <div style={{height: '100%'}}>
            <div className={styles['comment-main-div']} onScroll={handleScroll} >
              {comment.map((item, index)=> generateParentC(item, index))}
              {displaynomore ? <div className={styles['divider']}>没有更多评论</div> : <></>}
            </div>       
            <Tooltip position='top' trigger='hover' content={ t['comment.input.enter'] }>
              <TextArea
                autoComplete={'off'}
                maxLength={150}
                showWordLimit={true}
                className={styles['comment-input']}
                placeholder={ isLogin ? t['comment.input.placeholder'] :   t['comment.input.placeholder.plslog']}
                value = {valuebottom}
                onChange={(e)=>{
                  setValueBottom(e)
                }}
                disabled={!isLogin}
                onKeyDown={(e) => (handleKeyDownBottom(e, videoinfo['video_uid']))}
              />
            </Tooltip>
            
          </div>
      );
  }
  
  export default CommentDrawer;