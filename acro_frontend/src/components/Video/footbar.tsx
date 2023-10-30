import React, { forwardRef, useState } from 'react';
import { Button, Tooltip, Space, Slider, Switch } from '@arco-design/web-react';
import styles from './style/index.module.less';
import cs from 'classnames';
import { IconSound, IconMute } from '@arco-design/web-react/icon';
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
  const { icon, tooltip, onClick } = props;
  return (
    <Tooltip position='lt' trigger='hover' content={tooltip}>
        <Button
            icon={icon}
            shape="square"
            type="secondary"
            onClick={onClick}
            className={cs(styles['icon-foot-button'])}
        />
    </Tooltip>
  );
}

function FootBar(props, ref) {
  const { visible, playstate, timestate, playclick, volume, volumechange, autostate, setauto } = props;
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
                            tooltip={t['tooltip.like']}
                        />
                        <p className={styles['foot-time']}>{secondsToTimeFormat(timestate['now'])} {'/'} {secondsToTimeFormat(timestate['whole'])} </p>
                        <VolumeSlider value={volume} onChange={volumechange} />
                    </div>
                    <div className={styles['foot-group-right']}>
                        <Switch checkedText='Auto' uncheckedText='Auto' onChange={setauto} defaultChecked={autostate}/>
                    </div>
            </div>
            ) : <></>
        }
    </div>
  );
}

export default forwardRef(FootBar);
