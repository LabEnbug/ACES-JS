import React, { forwardRef, useState } from 'react';
import { Button, Tooltip, Space, Slider, Switch, Select, } from '@arco-design/web-react';
import styles from './style/index.module.less';
import cs from 'classnames';
import { IconSound, IconMute } from '@arco-design/web-react/icon';
import { FullScreen, OffScreen } from '@icon-park/react'
import { PlayOne, Pause } from '@icon-park/react'
import locale from './locale';
import useLocale from '@/utils/useLocale';

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
  const { visible, playstate, timestate, playclick, volume, volumechange, autostate, setauto, playbackrate, setplaybackrate, fullscreen, fullscreenchange } = props;
  const t = useLocale(locale);

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
