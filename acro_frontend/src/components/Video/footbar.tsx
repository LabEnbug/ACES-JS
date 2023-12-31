import React, { forwardRef, useState, useEffect } from 'react';
import { Button, Tooltip, Space, Slider, Switch, Select, Input, Avatar} from '@arco-design/web-react';
import styles from './style/index.module.less';
import cs from 'classnames';
import { IconSound, IconMute, IconUser } from '@arco-design/web-react/icon';
import { FullScreen, OffScreen } from '@icon-park/react'
import { PlayOne, Pause, Check, CheckOne, MoreTwo } from '@icon-park/react'
import locale from './locale';
import useLocale from '@/utils/useLocale';
import { useSelector } from 'react-redux';
import { GlobalState } from '@/store';

function secondsToTimeFormat(seconds) {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds =  Math.floor(seconds % 60);
    return `${minutes.toString().padStart(2, '0')}:${remainingSeconds.toString().padStart(2, '0')}`;
}

function VolumeSlider(props) {
    const {value, onChange} = props;
    return (
        <Space style={{'marginLeft': '5px'}}>
            <IconMute
            style={{
                fontSize: 20,
                color: value > 0 ? 'var(--color-text-4)' : 'var(--color-text-1)',
            }}
            />
            <Slider value={value} onChange={onChange} style={{ width: 150 }} />
            <IconSound
            style={{
                fontSize: 20,
                color: value === 0 ? 'var(--color-text-4)' : 'var(--color-text-1)',
            }}
            />
      </Space>
    )
}


function IconButton(props) {
  const { icon, tooltip, onClick, className } = props;
  return (
    <Tooltip position='lt' trigger='hover' content={tooltip}>
        <Button
            icon={icon}
            shape="square"
            type="secondary"
            onClick={onClick}
            className={cs(className)}
        />
    </Tooltip>
  );
}

function FootBar(props, ref) {
  const { visible, 
          playstate, 
          timestate, 
          playclick, 
          volume, 
          volumechange, 
          autostate, 
          setauto, 
          playbackrate, 
          setplaybackrate, 
          fullscreen, 
          fullscreenchange,
          sendBullet,
          closeBullet,
          bulletState,
        } = props;
  const t = useLocale(locale);
  const [bullet, setBullet] = useState('');
  const { isLogin } = useSelector((state: GlobalState) => state);
  return (
    <div>
        {
            visible ? (
                <div>
                    <div className={styles['foot-group-left']}>
                        <IconButton
                            icon={
                            <>
                            {
                                playstate ? <PlayOne theme="filled" size="36" fill="#ffffff"/> : <Pause theme="filled" size="36" fill="#ffffff"/>
                            }
                            </>
                            }
                            onClick = {playclick}
                            className = {styles['icon-foot-button']}
                        />
                        <p className={styles['foot-time']}>{secondsToTimeFormat(timestate['now'])} {'/'} {secondsToTimeFormat(timestate['whole'])} </p>
                        <VolumeSlider value={volume} onChange={volumechange} />
                        <div style={{marginLeft: '5px'}}>
                            <Input
                                autoComplete={'off'}
                                maxLength={20}
                                addBefore={(
                                    <Tooltip position='top' trigger='hover' content={ bulletState ? t['tooltip.bullets.open'] : t['tooltip.bullets.close'] }>
                                        <Avatar
                                            size={24}
                                            shape='square'
                                            triggerIcon={ bulletState ? <MoreTwo theme="filled" size="8" fill="#000000"/> : <Check theme="filled" size="8" fill="#ff2c55"/> }
                                            triggerIconStyle={{
                                                background: 'transparent'
                                            }}
                                            autoFixFontSize={true}
                                            onClick={closeBullet}
                                            style={{
                                                backgroundColor: '#ffffff',
                                            }}>
                                            <span style={{color:'#000000'}}>弹</span>
                                        </Avatar>
                                    </Tooltip>
                                    )}
                                placeholder={isLogin ? t['comment.input.placeholder'] :   t['comment.input.placeholder.plslog']}
                                value={bullet}
                                onChange={(e)=>{setBullet(e)}}
                                onPressEnter={(e)=>{
                                    sendBullet(e).then(res=>{
                                        setBullet('');
                                    }).catch((err)=>{
                                        console.error(err);
                                    });
                                }}
                                disabled={!isLogin}
                            />
                        </div>
                        {/* <Button onClick={handleSend}>发送</Button> */}
                    </div>
                    <div className={styles['foot-group-right']}>
                        <Switch checkedText={t['footbar.auto']} uncheckedText={t['footbar.auto']} onChange={setauto} defaultChecked={autostate} className={autostate ? styles['foot-autoplay-on'] : styles['foot-autoplay-off']}/>
                        <Select
                            triggerElement={<p className={styles['foot-playback']}>{'倍速'}</p>}
                            options={[
                                { label: '2.0X', value: 2 },
                                { label: '1.5X', value: 1.5 },
                                { label: '1.25X', value: 1.25 },
                                { label: '1.0X', value: 1 },
                                { label: '0.5X', value: 0.5 },
                            ]}
                            value={playbackrate}
                            triggerProps={{
                                autoAlignPopupWidth: false,
                                autoAlignPopupMinWidth: true,
                                position: 'top',
                            }}
                            trigger="hover"
                            onChange={setplaybackrate} />
                        <IconButton
                            icon={
                            <>
                            {
                                fullscreen ? <OffScreen theme="filled" size="28" fill="#ffffff"/> : <FullScreen theme="filled" size="28" fill="#ffffff"/>
                            }
                            </>
                            }
                            onClick = {fullscreenchange}
                            tooltip={t['tooltip.fullscreen']}
                            className = {styles['icon-foot-fullscreen-button']}
                        />
                    </div>
            </div>
            ) : <></>
        }
    </div>
  );
}

export default forwardRef(FootBar);
