import { useEffect, useState } from 'react';
import {  Tabs, Typography, Comment, Avatar, Input, Tooltip, Message } from '@arco-design/web-react';
import useLocale from '@/utils/useLocale';
import locale from './locale';
import { IconHeart, IconMessage, IconStar } from '@arco-design/web-react/icon';
import styles from './style/index.module.less';
import GetAxios from '@/utils/getaxios';
import GetUserInfo from "@/utils/getuserinfo";

const TextArea = Input.TextArea;
const TabPane = Tabs.TabPane;
const style = {
  textAlign: 'center',
  marginTop: 20,
  textAlign: 'left',
};


function CommentDrawer(props) {
    const {videoinfo} = props;
    const t = useLocale(locale);
    const [valuebottom, setValueBottom] = useState('');
    const [ comment, SetComment ] = useState([]);  
    const [ log, setLog ] = useState(false);
    const actions = (
        <span className='custom-comment-action'>
          <IconMessage /> Reply
        </span>
    );

    const JudgeStatus = (data: any) => {
      if (data.status != 200) {
        Message.error(data.err_msg);
        return false;
      }
      return true;
    }

    const handleKeyDownBottom = (e, uid) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        const baxio = GetAxios();
        const param = new FormData();
        // 如果只按下了Enter，阻止默认行为并触发你的事件
        e.preventDefault();
        if (e.target.value === '') {
          Message.error(t['comment.input.enter.empty']);
          return;
        }
        param.append('video_uid', uid);
        param.append('content', e.target.value);
        param.append('quote_comment_id', 0);
        baxio.post('v1-api/v1/video/comment/make', param).then(res=> {
          if (JudgeStatus(res)) {
            Message.info(t['comment.input.post.success']);
            setValueBottom('');
          }
        }).catch(e => {
          Message.info(t['comment.input.post.failed']);
          console.error(e);
        })
      }
    };

    useEffect(()=> {
      setLog(GetUserInfo() ? true : false);
    }, [log]);

    useEffect(()=> {
      const baxio = GetAxios();
      if (comment.length == 0) {

      }
    }, [comment]);
    

    return (
        <Tabs defaultActiveTab='1'>
          <TabPane key='1' title={t['comment']} style={{'color': '#ffffff'}}>
            <Typography.Paragraph style={style}>
            <Comment
                actions={actions}
                author={'Socrates'}
                avatar='//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/e278888093bef8910e829486fb45dd69.png~tplv-uwbnlip3yd-webp.webp'
                content={<div>Comment body content.</div>}
                datetime='1 hour'
                >
                <Comment
                    actions={actions}
                    author='Balzac'
                    avatar='//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/9eeb1800d9b78349b24682c3518ac4a3.png~tplv-uwbnlip3yd-webp.webp'
                    content={<div>Comment body content.</div>}
                    datetime='1 hour'
                >
                    <Comment
                    actions={actions}
                    author='Austen'
                    avatar='//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/8361eeb82904210b4f55fab888fe8416.png~tplv-uwbnlip3yd-webp.webp'
                    content={<div> Reply content </div>}
                    datetime='1 hour'
                    />
                    <Comment
                    actions={actions}
                    author='Plato'
                    avatar='//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/3ee5f13fb09879ecb5185e440cef6eb9.png~tplv-uwbnlip3yd-webp.webp'
                    content={<div> Reply content </div>}
                    datetime='1 hour'
                    />
                </Comment>
            </Comment>
            {
              log ?           
              <Tooltip position='top' trigger='hover' content={ t['comment.input.enter'] }>
                <TextArea
                  className={styles['comment-input']}
                  placeholder={t['comment.input.placeholder']}
                  // autoSize={{ minRows: 1, maxRows: 3 }}
                  // searchButton='Search'
                  // maxLength={120}
                  // showWordLimit
                  value = {valuebottom}
                  onChange={(e)=>{
                    setValueBottom(e)
                  }}
                  onKeyDown={(e) => (handleKeyDownBottom(e, videoinfo['video_uid']))}
                />
              </Tooltip> : <></> 
            
            }
            </Typography.Paragraph>
          </TabPane>
          <TabPane key='2' title='Tab 2' disabled>
            <Typography.Paragraph style={style}>Content of Tab Panel 2</Typography.Paragraph>
          </TabPane>
          <TabPane key='3' title='Tab 3'>
            <Typography.Paragraph style={style}>Content of Tab Panel 3</Typography.Paragraph>
          </TabPane>
        </Tabs>
      );
  }
  
  export default CommentDrawer;