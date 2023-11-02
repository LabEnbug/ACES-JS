import { useState } from 'react';
import {  Tabs, Typography, Comment, Avatar } from '@arco-design/web-react';
import useLocale from '@/utils/useLocale';
import locale from './locale';
import { IconHeart, IconMessage, IconStar } from '@arco-design/web-react/icon';

const TabPane = Tabs.TabPane;
const style = {
  textAlign: 'center',
  marginTop: 20,
  textAlign: 'left',
};

function CommentDrawer(props) {
    const t = useLocale(locale);
    const actions = (
        <span className='custom-comment-action'>
          <IconMessage /> Reply
        </span>
      );
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